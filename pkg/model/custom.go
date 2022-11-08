package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type StringMap map[string]string

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

type NullableString string

const NullString NullableString = "\x00"

// implements driver.Valuer, will be invoked automatically when written to the db
func (s NullableString) Value() (driver.Value, error) {
	if s == NullString {
		return nil, nil
	}
	return []byte(s), nil
}

// implements sql.Scanner, will be invoked automatically when read from the db
func (s *NullableString) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		*s = NullableString(v)
	case []byte:
		*s = NullableString(v)
	case nil:
		*s = ""
	}
	return nil
}
