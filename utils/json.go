package utils

import (
	"encoding/json"

	"github.com/TylerBrock/colorjson"
)

func colorJson(jsonData []byte) ([]byte, error) {
	var colorJsonMap interface{}
	err := json.Unmarshal(jsonData, &colorJsonMap)
	if err != nil {
		return nil, err
	}

	return colorjson.Marshal(colorJsonMap)
}