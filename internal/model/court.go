package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	return string(b), err
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}
	switch v := value.(type) {
	case string:
		if v == "" || v == "[]" {
			*s = StringArray{}
			return nil
		}
		return json.Unmarshal([]byte(v), s)
	case []byte:
		if len(v) == 0 || string(v) == "[]" {
			*s = StringArray{}
			return nil
		}
		return json.Unmarshal(v, s)
	default:
		return fmt.Errorf("cannot scan %T into StringArray", value)
	}
}

type Court struct {
	ID            uint        `gorm:"primaryKey"`
	Name          string
	Description   string      `gorm:"type:varchar(255);null"`
	PricePerHour  int
	Image         string
	Features      StringArray `gorm:"type:json"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}