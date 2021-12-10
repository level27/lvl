package utils

import (
	"encoding/json"

	"github.com/TylerBrock/colorjson"
)

func colorJson(jsonData []byte) ([]byte, error) {
	var colorJsonMap map[string]interface{}
	json.Unmarshal(jsonData, &colorJsonMap)
	return colorjson.Marshal(colorJsonMap)
}