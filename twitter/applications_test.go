package twitter

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestApplicationService_RateLimitStatus(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/1.1/application/rate_limit_status.json", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		assertQuery(t, map[string]string{"resources": "users,statuses"}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w,
			// https://dev.twitter.com/rest/reference/get/application/rate_limit_status
			`
			{
				"rate_limit_context": {
					"access_token": "786491-24zE39NUezJ8UTmOGOtLhgyLgCkPyY4dAcx6NA6sDKw"
				},
				"resources": {
					"users": {
						"/users/search": {
							"limit": 180,
							"remaining": 180,
							"reset": 1403602426
						}
					},
					"statuses": {
						"/statuses/mentions_timeline": {
							"limit": 15,
							"remaining": 15,
							"reset": 1403602426
						}
					}
				}
			}`,
		)
	})

	client := NewClient(httpClient)
	rateLimits, _, err := client.Applications.RateLimitStatus(&RateLimitStatusParams{Resources: []string{"users", "statuses"}})
	assert.Nil(t, err)

	_, ok := rateLimits.Resources["statuses"]
	require.True(t, ok, `expected rateLimits.Resources["statuses"]`)
	_, ok = rateLimits.Resources["statuses"]["/statuses/mentions_timeline"]
	require.True(t, ok, `expected rateLimits.Resources["statuses"]["/statuses/mentions_timeline"]`)
	assert.Equal(t, 15, rateLimits.Resources["statuses"]["/statuses/mentions_timeline"].Limit)

	_, ok = rateLimits.Resources["users"]
	require.True(t, ok, `expected rateLimits.Resources["users"]`)
	_, ok = rateLimits.Resources["users"]["/users/search"]
	require.True(t, ok, `expected rateLimits.Resources["users"]["/users/search"]`)
	assert.Equal(t, EpochTime(time.Unix(1403602426, 0)), rateLimits.Resources["users"]["/users/search"].Reset)
}

func TestApplicationService_RateLimitStatusNilParams(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/1.1/application/rate_limit_status.json", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		assertQuery(t, map[string]string{}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w,
			// https://dev.twitter.com/rest/reference/get/application/rate_limit_status
			`
			{
				"rate_limit_context": {
					"access_token": "786491-24zE39NUezJ8UTmOGOtLhgyLgCkPyY4dAcx6NA6sDKw"
				},
				"resources": {
					"users": {
						"/users/search": {
							"limit": 180,
							"remaining": 180,
							"reset": 1403602426
						}
					},
					"statuses": {
						"/statuses/mentions_timeline": {
							"limit": 15,
							"remaining": 15,
							"reset": 1403602426
						}
					}
				}
			}`,
		)
	})

	client := NewClient(httpClient)
	rateLimits, _, err := client.Applications.RateLimitStatus(nil)
	assert.Nil(t, err)

	_, ok := rateLimits.Resources["statuses"]
	require.True(t, ok, `expected rateLimits.Resources["statuses"]`)
	_, ok = rateLimits.Resources["statuses"]["/statuses/mentions_timeline"]
	require.True(t, ok, `expected rateLimits.Resources["statuses"]["/statuses/mentions_timeline"]`)
	assert.Equal(t, 15, rateLimits.Resources["statuses"]["/statuses/mentions_timeline"].Limit)

	_, ok = rateLimits.Resources["users"]
	require.True(t, ok, `expected rateLimits.Resources["users"]`)
	_, ok = rateLimits.Resources["users"]["/users/search"]
	require.True(t, ok, `expected rateLimits.Resources["users"]["/users/search"]`)
	assert.Equal(t, EpochTime(time.Unix(1403602426, 0)), rateLimits.Resources["users"]["/users/search"].Reset)
}

func TestEpochTime(t *testing.T) {
	tm := EpochTime(time.Unix(1403602426, 0))
	jsonStr, err := json.Marshal(&tm)
	assert.NoError(t, err)
	assert.Equal(t, "1403602426", string(jsonStr))

	v := &RateLimit{}
	err = json.Unmarshal([]byte(`{
		"limit": 15,
		"remaining": 15,
		"reset": "xasdasa"
	}`), v)

	assert.EqualError(t, err, "json: cannot unmarshal string into Go value of type int64")

	var p *EpochTime
	// nil marshal
	jsonStr, err = json.Marshal(p)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(jsonStr))
}
