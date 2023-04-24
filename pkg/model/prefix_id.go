package model

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"strings"
)

type UserId string
type PlatformId string
type NetworkId string
type AssetId string
type DeviceId string
type ContactId string
type LocationId string
type InstrumentId string
type TxLegId string
type TransactionId string
type AuthStrategyId string
type ApiKeyId string
type ContractId string

func (id UserId) Value() (driver.Value, error)         { return value(id) }
func (id PlatformId) Value() (driver.Value, error)     { return value(id) }
func (id NetworkId) Value() (driver.Value, error)      { return value(id) }
func (id AssetId) Value() (driver.Value, error)        { return value(id) }
func (id DeviceId) Value() (driver.Value, error)       { return value(id) }
func (id ContactId) Value() (driver.Value, error)      { return value(id) }
func (id LocationId) Value() (driver.Value, error)     { return value(id) }
func (id InstrumentId) Value() (driver.Value, error)   { return value(id) }
func (id TxLegId) Value() (driver.Value, error)        { return value(id) }
func (id TransactionId) Value() (driver.Value, error)  { return value(id) }
func (id AuthStrategyId) Value() (driver.Value, error) { return value(id) }
func (id ApiKeyId) Value() (driver.Value, error)       { return value(id) }
func (id ContractId) Value() (driver.Value, error)     { return value(id) }

func (id UserId) Scan(src interface{}) error         { return scan(id, src) }
func (id PlatformId) Scan(src interface{}) error     { return scan(id, src) }
func (id NetworkId) Scan(src interface{}) error      { return scan(id, src) }
func (id AssetId) Scan(src interface{}) error        { return scan(id, src) }
func (id DeviceId) Scan(src interface{}) error       { return scan(id, src) }
func (id ContactId) Scan(src interface{}) error      { return scan(id, src) }
func (id LocationId) Scan(src interface{}) error     { return scan(id, src) }
func (id InstrumentId) Scan(src interface{}) error   { return scan(id, src) }
func (id TxLegId) Scan(src interface{}) error        { return scan(id, src) }
func (id TransactionId) Scan(src interface{}) error  { return scan(id, src) }
func (id AuthStrategyId) Scan(src interface{}) error { return scan(id, src) }
func (id ApiKeyId) Scan(src interface{}) error       { return scan(id, src) }
func (id ContractId) Scan(src interface{}) error     { return scan(id, src) }

// func Trim(id string) string {
// 	parts := strings.Split(id, "_")
// 	if len(parts) < 2 {
// 		return ""
// 	}
// 	return parts[1]
// }

func Trim(id interface{}) string {
	str, ok := id.(string)
	if !ok {
		return ""
	}
	parts := strings.Split(str, "_")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

var typeToPrefix = map[string]string{
	"UserId":         "user",
	"PlatformId":     "platform",
	"NetworkId":      "network",
	"AssetId":        "asset",
	"DeviceId":       "device",
	"ContactId":      "contact",
	"LocationId":     "location",
	"InstrumentId":   "instrument",
	"TxLegId":        "txleg",
	"TransactionId":  "tx",
	"AuthStrategyId": "auth",
	"ApiKeyId":       "apikey",
	"ContractId":     "contract",
}

func value(id interface{}) (driver.Value, error) {
	t := reflect.TypeOf(id).Name()
	prefix, ok := typeToPrefix[t]
	if !ok {
		return "", errors.New("unknown type")
	}
	if id.(string)[0:len(prefix)+1] != prefix+"_" {
		return nil, errors.New("invalid prefix")
	}
	return []byte(id.(string)[len(prefix)+1:]), nil
}

func scan(destination interface{}, src interface{}) error {
	t := reflect.TypeOf(destination).Name()
	prefix, ok := typeToPrefix[t]
	if !ok {
		return errors.New("unknown type")
	}
	switch v := src.(type) {
	case string:
		destination = prefix + "_" + v
	case []byte:
		destination = prefix + "_" + string(v)
	case nil:
		destination = ""
	}
	return nil
}

////////

// type ID struct {
// 	Prefix string
// 	Id     string
// }

// func (id ID) ToString() string {
// 	return id.Id
// }

// func (id *ID) Scan(src interface{}) error {
// 	switch v := src.(type) {
// 	case string:
// 		id.Id = v
// 	case []byte:
// 		id.Id = string(v)
// 	case nil:
// 		id.Id = ""
// 	}
// 	return nil
// }

// func (id ID) Value() (driver.Value, error) {
// 	return []byte(id.Id), nil
// }

// // Marshal the ID
// func (id ID) MarshalJSON() ([]byte, error) {
// 	return []byte(id.Id), nil
// }

// // Unmarshal the ID
// func (id *ID) UnmarshalJSON(data []byte) error {
// 	id.Id = string(data)
// 	return nil
// }

// type UserId struct{ ID }
// type PlatformId struct{ ID }
// type NetworkId struct{ ID }
// type AssetId struct{ ID }
// type DeviceId struct{ ID }
// type ContactId struct{ ID }
// type LocationId struct{ ID }
// type InstrumentId struct{ ID }
// type TxLegId struct{ ID }
// type TransactionId struct{ ID }
// type AuthStrategyId struct{ ID }
// type ApiKeyId struct{ ID }
// type ContractId struct{ ID }
