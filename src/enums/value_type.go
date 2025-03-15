package enums

import "database/sql/driver"

type ValueType string

const (
	Int      ValueType = "int"
	Float    ValueType = "float"
	String   ValueType = "string"
	Bool     ValueType = "bool"
	Date     ValueType = "date"
	DateTime ValueType = "datetime"
	JSON     ValueType = "json"
)

func (vt *ValueType) Scan(value interface{}) error {
	*vt = ValueType(value.(string))
	return nil
}

func (vt ValueType) Value() (driver.Value, error) {
	return string(vt), nil
}

func (vt ValueType) String() string {
	return string(vt)
}
