package ratelimit

import "time"

var MaxRequestsPerMinute = 500

// newExpiration returns a new time.Time adding the passed seconds to the current time
func newExpiration(seconds int) time.Time {

	return time.Now().Add(time.Duration(seconds) * time.Second)
}

// RateLimit contains the amount of request already made and the expiration time
type RateLimit struct {
	amount     int
	expiration time.Time
}

// NewRateLimit returns a new *RateLimit with a default amount of 1 and a new expiration time of the current time + 60 seconds
func NewRateLimit() *RateLimit { return &RateLimit{amount: 1, expiration: newExpiration(60)} }

// Expired returns if the *RateLimit is expired or not
func (r *RateLimit) Expired() bool { return time.Now().Before(time.Now()) }

func (r *RateLimit) Reset() { r = NewRateLimit() }

func (r *RateLimit) AddAmount() { r.amount++ }

func (r *RateLimit) Limited() bool { return r.amount >= MaxRequestsPerMinute }

func (r *RateLimit) Expiration() time.Time { return r.expiration }

func (r *RateLimit) Amount() int { return r.amount }
