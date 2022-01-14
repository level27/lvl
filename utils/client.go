package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"bitbucket.org/level27/lvl/types"
	"github.com/Jeffail/gabs/v2"
)

var TraceRequests bool

// Client defines the API Client structure
type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

// NewAPIClient creates a client for doing the API calls
func NewAPIClient(uri string, apiKey string) *Client {
	return &Client{
		BaseURL: uri,
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  struct {
		Children struct {
			Content struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"content,omitempty"`
			SSLForce struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"sslForce,omitempty"`
			SSLCertificate struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"sslCertificate,omitempty"`
			HandleDNS struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"handleDns,omitempty"`
			Authentication struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"authentication,omitempty"`
			Appcomponent struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"appcomponent,omitempty"`
		} `json:"children"`
	} `json:"errors"`
}

func (er errorResponse) Error() string {
	var sb strings.Builder
	sb.WriteString(er.Message)

	fields := reflect.TypeOf(er.Errors.Children)
	values := reflect.ValueOf(er.Errors.Children)

	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		if value.Field(0).Len() > 0 {
			sb.WriteString(fmt.Sprintf("\n%v = %v", field.Name, value))
		}
	}

	return sb.String()
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func (c *Client) sendRequestRaw(method string, endpoint string, data interface{}, headers map[string]string) (*http.Response, error) {
	reqData := bytes.NewBuffer([]byte(nil))
	if data != nil {
		str, ok := data.(string)
		if ok {
			reqData = bytes.NewBuffer([]byte(str))
		} else {
			jsonDat, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			reqData = bytes.NewBuffer(jsonDat)
		}
	}

	fullUrl := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)

	if TraceRequests {
		fmt.Fprintf(os.Stderr, "Request: %s %s\n", method, fullUrl)
		if reqData.Len() != 0 {
			colored, err := colorJson(reqData.Bytes())
			var str string
			if err == nil {
				str = string(colored)
			} else {
				str = string(reqData.String())
			}

			fmt.Fprintf(os.Stderr, "Request Body: %s\n", str)
		}
	}

	req, err := http.NewRequest(method, fullUrl, reqData)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", "level27_lvl/1.0")
	req.Header.Set("Authorization", c.apiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if TraceRequests {
		fmt.Fprintf(os.Stderr, "Response: %d %s\n", res.StatusCode, http.StatusText(res.StatusCode))
	}

	return res, err
}

func (c *Client) sendRequest(method string, endpoint string, data interface{}) ([]byte, error) {
	headers := map[string]string{"Accept": "application/json"}
	if data != nil {
		headers["Content-Type"] = "application/json"
	}

	res, err := c.sendRequestRaw(method, endpoint, data, headers)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if method == "UPDATE" && res.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if TraceRequests && len(body) != 0 {
		bodyPrint := body
		if json.Valid(bodyPrint) {
			bodyPrint, _ = colorJson(bodyPrint)
		}
		fmt.Fprintf(os.Stderr, "Response Body: %s\n", string(bodyPrint))
	}

	if isErrorCode(res.StatusCode) {
		return nil, formatRequestError(res.StatusCode, body)
	}

	return body, nil
}

func isErrorCode(statusCode int) bool {
	return statusCode < http.StatusOK || statusCode >= http.StatusBadRequest
}

func formatRequestError(statusCode int, body []byte) error {
	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		return err
	}

	// log.Printf("client.go: ERROR: %v", jsonParsed)
	for key, child := range jsonParsed.Search("errors", "children").ChildrenMap() {
		if child.Data().(map[string]interface{})["errors"] != nil {
			errorMessages := child.Data().(map[string]interface{})["errors"].([]interface{})
			if len(errorMessages) > 0 {
				for _, err := range errorMessages {
					log.Printf("Key=>%v, Value=>%v\n", key, err)
					return fmt.Errorf("%v : %v", key, err)
				}
			}
		}
	}

	var errRes errorResponse
	if err = json.Unmarshal(body, &errRes); err == nil {
		return errRes
	}

	return fmt.Errorf("unknown error, status code: %d", statusCode)
}

func (c *Client) invokeAPI(method string, endpoint string, data interface{}, result interface{}) error {
	body, err := c.sendRequest(method, endpoint, data)

	if err != nil {
		return err
	}

	if result != nil {

		err = json.Unmarshal(body, &result)
	}

	return err
}

func AssertApiError(e error, directory string) {
	TranslateStatusCode(e, directory)
	if e != nil {

		log.Fatalf("client.go: API error - %s\n", e.Error())
	}
}

func TranslateStatusCode(e error, directory string) {
	if e != nil {
		splittedError := strings.Split(e.Error(), " ")
		var result string
		status := splittedError[len(splittedError)-1]
		switch status {
		case "204":
			result = fmt.Sprintf("Status: %v. Request succesfully executed", status)
		case "400":
			result = fmt.Sprintf("Status: %v. Bad request", status)
		case "403":
			result = fmt.Sprintf("Status: %v. You do not have acces to this %v", status, directory)
		case "404":
			result = fmt.Sprintf("Status: %v. %v not found", status, directory)
		case "500":
			result = fmt.Sprintf("Status: %v. You have no proper rights to acces the controller", status)
		default:
			result = "No Status code received"
		}

		log.Println(result)
	} else {
		log.Println("Request succesfully executed")
	}


}

func formatCommonGetParams(params types.CommonGetParams) string {
	return fmt.Sprintf("limit=%d&filter=%s", params.Limit, url.QueryEscape(params.Filter))
}
