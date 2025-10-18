package server

import (
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var global_identifier = "global"

type RateLimiter struct {
	RequestsAllowedInWindow int           // number of requests allowed in the sliding time window
	WindowSize              time.Duration // duration of the sliding rate limit window
	UsingIP                 bool          // rate limit by IP if true, otherwise global rate limit

	requests map[string][]time.Time // map of IP/or "global" to slice of request timestamps
	mx       sync.Mutex
}

func (r *RateLimiter) AllowRequest(identifier string) bool {
	r.mx.Lock()
	defer r.mx.Unlock()

	now := time.Now()

	idReqs := r.cleanupIdentifier(now, identifier)
	if len(idReqs) > r.RequestsAllowedInWindow { /// immediately return false if over limit
		return false
	}

	r.requests[identifier] = append(r.requests[identifier], now)
	return true
}

// cleanupIdentifier removes timestamps outside the sliding window and returns
// the current count of requests within the window.
func (r *RateLimiter) cleanupIdentifier(now time.Time, identifier string) []time.Time {
	identifierRequests := r.requests[identifier]

	// as long as idreqs has entries and the time between now and oldest entry is
	// greater than window size, remove oldest entry
	// if the identifier isn't in the map, this loop won't run bc len(nil map) == 0
	for len(identifierRequests) > 0 && now.Sub(identifierRequests[0]) > r.WindowSize {
		identifierRequests = identifierRequests[1:]
	}
	r.requests[identifier] = identifierRequests
	return identifierRequests
}

func (r *RateLimiter) cleanup() {
	r.mx.Lock()
	defer r.mx.Unlock()

	now := time.Now()
	for identifier := range r.requests {
		if len(r.cleanupIdentifier(now, identifier)) == 0 {
			delete(r.requests, identifier) // delete if no requests remain (no requests in the last window)
		}
	}
}

func (r *RateLimiter) Start() *RateLimiter {
	// cleanup every window size interval
	go func() {
		ticker := time.NewTicker(r.WindowSize)
		defer ticker.Stop() // this goroutine runs forever, but good practice in case of future changes
		for range ticker.C {
			r.cleanup()
		}
	}()
	return r
}

func (r *RateLimiter) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		identifier := global_identifier
		if r.UsingIP {
			identifier = c.RealIP()
		}
		if !r.AllowRequest(identifier) {
			return c.String(429, "too many requests")
		}
		return next(c)
	}
}

func NewRateLimiter(requestsAllowedInWindow int, windowSize time.Duration, usingIP bool) *RateLimiter {
	return &RateLimiter{
		RequestsAllowedInWindow: requestsAllowedInWindow,
		WindowSize:              windowSize,
		UsingIP:                 usingIP,
		requests:                make(map[string][]time.Time),
	}
}
