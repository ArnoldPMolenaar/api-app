package enums

import "database/sql/driver"

type Level string

const (
	Public  Level = "public"
	Private Level = "private"
	Both    Level = "both"
)

func (l *Level) Scan(value interface{}) error {
	*l = Level(value.(string))
	return nil
}

func (l Level) Value() (driver.Value, error) {
	return string(l), nil
}

func (l Level) String() string {
	return string(l)
}
