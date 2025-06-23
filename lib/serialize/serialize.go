package serialize

import (
	"encoding/json"
)

type TestData struct {
	Id   int
	Name string
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
