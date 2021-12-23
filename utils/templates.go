package utils

import (
	"fmt"
	"text/template"
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
	}
}