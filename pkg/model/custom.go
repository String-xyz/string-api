package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/lib/pq"
)

type StringMap map[string]string
type StringArray pq.StringArray

func (sm StringMap) Value() (driver.Value, error) {
	return json.Marshal(sm)
}

func (sm *StringMap) Scan(src interface{}) error {
	switch t := src.(type) {
	case string:
		return json.Unmarshal([]byte(t), sm)
	case []byte:
		return json.Unmarshal([]byte(t), sm)
	}
	return errors.New("unknown type")
}
