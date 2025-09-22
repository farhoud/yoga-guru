package utils

import (
	"encoding/json"
	"time"
)

// Define the custom time format.
const timeFormat = "15:04:05"

// CustomTime is a custom type to handle time-only JSON strings.
type CustomTime time.Time

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := time.Parse(timeFormat, s)
	if err != nil {
		return err
	}
	*ct = CustomTime(t)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	s := t.Format(timeFormat)
	return json.Marshal(s)
}
