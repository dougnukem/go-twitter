package twitter

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	// Twitter usage of Ruby time format
	tmStr := `"Mon Nov 29 21:18:15 +0000 2010"`
	tm := &Time{}
	err := json.Unmarshal([]byte(tmStr), tm)
	assert.NoError(t, err)

	expectedTime, err := time.Parse(time.RubyDate, "Mon Nov 29 21:18:15 +0000 2010")
	assert.NoError(t, err)

	assert.Equal(t, expectedTime, time.Time(*tm))
}

func TestTime_UnmarshalJSONError(t *testing.T) {
	tmStr := "asdfasdfasdf23423"
	tm := &Time{}
	err := json.Unmarshal([]byte(tmStr), tm)
	assert.EqualError(t, err, "invalid character 'a' looking for beginning of value")
}

func TestTime_MarshalJSON(t *testing.T) {
	tm := Time(time.Unix(1000, 0))
	tmJSON, err := json.Marshal(&tm)
	assert.NoError(t, err)
	assert.Equal(t, `"Wed Dec 31 18:16:40 -0600 1969"`, string(tmJSON))
}
