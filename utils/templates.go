package utils

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"
	"time"

	"github.com/fatih/color"
	"github.com/level27/l27-go"
)

// Returns the template helper functions accessible to our templates.
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
		"vt": func(colorCode string) string {
			if color.NoColor {
				return ""
			}
			return vtCodes[colorCode]
		},
	}
}

var vtCsi = "\x1B["
var vtCodes = map[string]string{
	"reset":         vtCsi + "0m",
	"black":         vtCsi + "30m",
	"red":           vtCsi + "31m",
	"green":         vtCsi + "32m",
	"yellow":        vtCsi + "33m",
	"blue":          vtCsi + "34m",
	"magenta":       vtCsi + "35m",
	"cyan":          vtCsi + "36m",
	"white":         vtCsi + "37m",
	"brightblack":   vtCsi + "90m",
	"brightred":     vtCsi + "91m",
	"brightgreen":   vtCsi + "92m",
	"brightyellow":  vtCsi + "93m",
	"brightblue":    vtCsi + "94m",
	"brightmagenta": vtCsi + "95m",
	"brightcyan":    vtCsi + "96m",
	"brightwhite":   vtCsi + "97m",
}

func jobChildCountTotalRecurse(job l27.Job) int {
	counter := len(job.Jobs)
	for _, j := range job.Jobs {
		counter += jobChildCountTotalRecurse(j)
	}
	return counter
}

func FormatUnixTimeF(seconds interface{}, fmt string) string {
	if seconds == nil {
		return "null"
	}
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
		if success {
			secs = int64(secsFloat)
		} else {
			secs = seconds.(int64)
		}
	}

	reqTime := time.Unix(secs, 0)
	return reqTime.Format(fmt)
}

// Format a unix time value returned by the API in a way that is human-readable.
func FormatUnixTime(seconds interface{}) string {
	result := FormatUnixTimeF(seconds, time.RFC1123)
	if result == "1970 Jan  1 01:00:00" || result == "Thu, 01 Jan 1970 01:00:00 CET" {
		return "null"
	} else {
		return result
	}
}
