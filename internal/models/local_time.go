package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const localTimeLayout = "2006-01-02 15:04:05"
const localDateLayout = "2006-01-02"

var localTimeLocation = time.FixedZone("Asia/Bangkok", 7*60*60)

type LocalTime time.Time

type LocalDate time.Time

func NewLocalTime(t time.Time) LocalTime {
	return LocalTime(t.In(localTimeLocation))
}

func NewLocalDate(t time.Time) LocalDate {
	tt := t.In(localTimeLocation)
	date := time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, localTimeLocation)
	return LocalDate(date)
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return []byte(`""`), nil
	}
	return json.Marshal(tt.In(localTimeLocation).Format(localTimeLayout))
}

func (t LocalDate) MarshalJSON() ([]byte, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return []byte(`""`), nil
	}
	return json.Marshal(tt.In(localTimeLocation).Format(localDateLayout))
}

func (t *LocalTime) UnmarshalJSON(b []byte) error {
	if t == nil {
		return errors.New("LocalTime: nil receiver")
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		*t = LocalTime(time.Time{})
		return nil
	}
	parsed, err := time.ParseInLocation(localTimeLayout, s, localTimeLocation)
	if err != nil {
		return err
	}
	*t = LocalTime(parsed)
	return nil
}

func (t *LocalDate) UnmarshalJSON(b []byte) error {
	if t == nil {
		return errors.New("LocalDate: nil receiver")
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		*t = LocalDate(time.Time{})
		return nil
	}
	parsed, err := time.ParseInLocation(localDateLayout, s, localTimeLocation)
	if err != nil {
		return err
	}
	*t = LocalDate(parsed)
	return nil
}

func (t LocalTime) Value() (driver.Value, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return nil, nil
	}
	return tt.In(localTimeLocation), nil
}

func (t LocalDate) Value() (driver.Value, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return nil, nil
	}
	date := time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, localTimeLocation)
	return date, nil
}

func (t *LocalTime) Scan(value interface{}) error {
	if t == nil {
		return errors.New("LocalTime: nil receiver")
	}
	switch v := value.(type) {
	case time.Time:
		*t = LocalTime(v.In(localTimeLocation))
		return nil
	case []byte:
		return t.parseString(string(v))
	case string:
		return t.parseString(v)
	case nil:
		*t = LocalTime(time.Time{})
		return nil
	default:
		return fmt.Errorf("LocalTime: unsupported type %T", value)
	}
}

func (t *LocalDate) Scan(value interface{}) error {
	if t == nil {
		return errors.New("LocalDate: nil receiver")
	}
	switch v := value.(type) {
	case time.Time:
		date := time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, localTimeLocation)
		*t = LocalDate(date)
		return nil
	case []byte:
		return t.parseString(string(v))
	case string:
		return t.parseString(v)
	case nil:
		*t = LocalDate(time.Time{})
		return nil
	default:
		return fmt.Errorf("LocalDate: unsupported type %T", value)
	}
}

func (t *LocalTime) parseString(s string) error {
	if s == "" {
		*t = LocalTime(time.Time{})
		return nil
	}
	parsed, err := time.ParseInLocation(localTimeLayout, s, localTimeLocation)
	if err != nil {
		return err
	}
	*t = LocalTime(parsed)
	return nil
}

func (t *LocalDate) parseString(s string) error {
	if s == "" {
		*t = LocalDate(time.Time{})
		return nil
	}
	parsed, err := time.ParseInLocation(localDateLayout, s, localTimeLocation)
	if err != nil {
		return err
	}
	*t = LocalDate(parsed)
	return nil
}
