package twitter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

// ApplicationService provides a method for application rate limits.
type ApplicationService struct {
	sling *sling.Sling
}

// newAccountService returns a new AccountService.
func newApplicationService(sling *sling.Sling) *ApplicationService {
	return &ApplicationService{
		sling: sling.Path("application/"),
	}
}

// EpochTime is an integer timestamp in seconds from UTC epoch time:
// https://en.wikipedia.org/wiki/Unix_time
type EpochTime time.Time

func (t EpochTime) String() string {
	return time.Time(t).String()
}

// UnmarshalJSON decodes a string using the time.Unix format
func (t *EpochTime) UnmarshalJSON(b []byte) error {
	var secs int64

	if err := json.Unmarshal(b, &secs); err != nil {
		return err
	}

	*t = EpochTime(time.Unix(secs, 0))
	return nil
}

// MarshalJSON encodes the time in the time.Unix format
func (t *EpochTime) MarshalJSON() ([]byte, error) {
	if t == nil {
		return nil, nil
	}
	timeStr := fmt.Sprintf(`%d`, time.Time(*t).Unix())
	return []byte(timeStr), nil
}

// RateLimit represents the rate limit status of a particular
// api endpoint
type RateLimit struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     EpochTime `json:"reset"`
}

// RateLimitStatus represents the rate limit status for
// all api endpoints
type RateLimitStatus struct {
	RateLimitContext struct {
		AccessToken string `json:"access_token"`
	} `json:"rate_limit_context"`
	Resources map[string]map[string]RateLimit `json:"resources"`
}

// RateLimitStatusParams are the params for ApplicationService.RateLimitStatus.
type RateLimitStatusParams struct {
	Resources []string `url:"resources,omitempty,comma"`
}

// RateLimitStatus returns the rate limits for this application (or user of this application)
// returns an error otherwise.
// Requires a user auth context.
// https://dev.twitter.com/rest/reference/get/application/rate_limit_status
func (s *ApplicationService) RateLimitStatus(params *RateLimitStatusParams) (*RateLimitStatus, *http.Response, error) {
	if params == nil {
		params = &RateLimitStatusParams{}
	}
	rateStatus := new(RateLimitStatus)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("rate_limit_status.json").QueryStruct(params).Receive(rateStatus, apiError)
	return rateStatus, resp, relevantError(err, *apiError)
}
