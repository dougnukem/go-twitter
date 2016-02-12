package twitter

import (
	"encoding/json"
	"fmt"
	"time"
)

// Time is a twitter specific time type. Twitter time
// fields are encoded using the time.RubyDate format
// so we have to custom json decode
type Time time.Time

// UnmarshalJSON decodes a string using the time.RubyDate format
func (t *Time) UnmarshalJSON(b []byte) error {
	var str string

	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	t2, err := time.Parse(time.RubyDate, str)

	if err != nil {
		return err
	}

	*t = Time(t2)
	return nil
}

// MarshalJSON encodes the time in the time.RubyDate format
func (t *Time) MarshalJSON() ([]byte, error) {
	if t == nil {
		return nil, nil
	}
	timeStr := fmt.Sprintf(`"%s"`, time.Time(*t).Format(time.RubyDate))
	return []byte(timeStr), nil
}
