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

// Convert a JSON-serializable value to a model of interface{} maps and slices. Effectively serializing the object to maps/slices with the JSON model.
func RoundTripJson(obj interface{}) interface{} {
	// Round-trip through JSON so we use the JSON (camelcased) key names in the YAML without having to re-define them
	bJson, _ := json.Marshal(obj)
	var interf interface{}
	json.Unmarshal(bJson, &interf)
	return interf
}
