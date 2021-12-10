package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

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
	s := fmt.Sprintf("%s\n", er.Message)
	fields := reflect.TypeOf(er.Errors.Children)
	values := reflect.ValueOf(er.Errors.Children)

	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		s += fmt.Sprintf("%v = %v\n", field.Name, value)
	}
	return s
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func (c *Client) sendRequest(method string, endpoint string, data interface{}) ([]byte, error) {

	reqData := bytes.NewBuffer([]byte(nil))
	if data != nil {
		jsonDat, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		reqData = bytes.NewBuffer(jsonDat)
	}

	fullUrl := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)

	if (TraceRequests) {
		fmt.Fprintf(os.Stderr, "Request: %s %s\n", method, fullUrl)
		if reqData.Len() != 0 {
			fmt.Fprintf(os.Stderr, "Request Body: %s\n", string(reqData.Bytes()))
		}
	}

	req, err := http.NewRequest(method, fullUrl, reqData)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "level27_lvl/1.0")
	req.Header.Set("Authorization", c.apiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if TraceRequests {
		fmt.Fprintf(os.Stderr, "Response: %d %s\n", res.StatusCode, http.StatusText(res.StatusCode))
	}

	if method == "UPDATE" && res.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if TraceRequests && len(body) != 0 {
		fmt.Fprintf(os.Stderr, "Response Body: %s\n", string(body))
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		jsonParsed, err := gabs.ParseJSON(body)
		if err != nil {
			return nil, err
		}

		// log.Printf("client.go: ERROR: %v", jsonParsed)
		for key, child := range jsonParsed.Search("errors", "children").ChildrenMap() {
			if child.Data().(map[string]interface{})["errors"] != nil {
				errorMessages := child.Data().(map[string]interface{})["errors"].([]interface{})
				if len(errorMessages) > 0 {
					for _, err := range errorMessages {
						log.Printf("Key=>%v, Value=>%v\n", key, err)
						return nil, fmt.Errorf("%v : %v", key, err)
					}
				}
			}
		}

		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return nil, errRes
		}
		
		return nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	return body, nil
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

func AssertApiError(e error) {
	if e != nil {
		log.Fatalf("client.go: API error - %s\n", e.Error())
	}
}
