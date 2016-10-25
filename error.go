package algorithmia

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

type Err struct {
	Message    string `json:"message" mapstructure:"message"`
	Stacktrace string `json:"stacktrace" mapstructure:"stacktrace"`
}

func (e *Err) Error() string {
	return e.Message
}

func ErrFromJsonData(data []byte) (*Err, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	if mm, ok := m["error"]; ok {
		var err Err
		if e := mapstructure.Decode(mm, &err); e != nil {
			return nil, e
		}
		return &err, nil
	}

	return nil, nil
}

func ErrorFromJsonData(data []byte) error {
	e1, e2 := ErrFromJsonData(data)
	if e2 != nil {
		return e2
	}
	if e1 != nil {
		return e1
	}
	return nil
}
