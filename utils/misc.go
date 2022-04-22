package utils

import (
	"mime"
	"net/http"
)

func parseContentDispositionFilename(resp *http.Response, fallback string) string {
	contentDisp := resp.Header.Get("Content-Disposition")
	if contentDisp != "" {
		_, params, err := mime.ParseMediaType(contentDisp)
		if err != nil {
			return fallback
		}

		return params["filename"]
	}

	return fallback
}