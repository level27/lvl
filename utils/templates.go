package utils

import (
	"fmt"
	"strconv"
	"text/template"
	"time"
)

func TemplateHelpers() template.FuncMap {
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
		"formatUnixTime": func(secs string) string {
			return FormatUnixTime(secs)
		},
	}
}


func FormatUnixTime(secondsString string) string {
	secs, err := strconv.ParseInt(secondsString, 10, 64)
	if err != nil {
		return ""
	}

	reqTime := time.Unix(secs, 0)
	return reqTime.Format("2006 Jan _2 15:04:05")
}
