package utils

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"
	"time"

	"bitbucket.org/level27/lvl/types"
)

func MakeTemplateHelpers(t *template.Template) template.FuncMap {
	return template.FuncMap{
		// Formats a currenct + price pair such as "EUR", "100" -> 1.00 €.
		"formatCurrency": func(currency string, price string) string {
			switch currency {
			case "EUR":
				beforeDecimal := price[:len(price)-2]
				afterDecimal := price[len(price)-2:]
				return fmt.Sprintf("%s.%s €", beforeDecimal, afterDecimal)
			}

			panic("Unknown currency: " + currency)
		},
		// Formats a string unix time and prints it
		"formatUnixTime": FormatUnixTime,
		// Formats a string unix time and prints it
		"formatUnixTimeF": FormatUnixTimeF,
		"include": func(name string, data interface{}) (string, error) {
			// https://stackoverflow.com/a/59401696
			buf := bytes.NewBuffer(nil)
			if t.ExecuteTemplate(buf, name, data) != nil {
				return "", nil
			}

			return buf.String(), nil
		},
		"jobStatusSafe": func(status interface{}) (int, error) {
			str, success := status.(string)
			if success {
				if str == "busy" {
					return 999, nil
				}

				return 0, fmt.Errorf("unknown job status: %s", str)
			} else {
				i, success := status.(int)
				if !success {
					return 0, fmt.Errorf("unknown job status type")
				}

				return i, nil
			}
		},
		"jobChildCountTotal": jobChildCountTotalRecurse,
	}
}

func jobChildCountTotalRecurse(job types.Job) int {
	counter := len(job.Jobs)
	for _, j := range job.Jobs {
		counter += jobChildCountTotalRecurse(j)
	}
	return counter
}

func FormatUnixTimeF(seconds interface{}, fmt string) string {
	var secs int64
	secString, success := seconds.(string)
	if success {
		secsParsed, err := strconv.ParseInt(secString, 10, 64)
		secs = secsParsed
		if err != nil {
			return ""
		}
	} else {
		secsFloat, success := seconds.(float64)
		if (success) {
			secs = int64(secsFloat)
		} else {
			secs = seconds.(int64)
		}
	}

	reqTime := time.Unix(secs, 0)
	return reqTime.Format(fmt)
}

func FormatUnixTime(seconds interface{}) string {
	return FormatUnixTimeF(seconds, "2006 Jan _2 15:04:05")
}