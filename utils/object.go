package utils

import "encoding/json"

// MarshalJSON marshals the object to JSON
func MarshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
