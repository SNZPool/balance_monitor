package balance_monitor

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

// EndpointLimiters rate-limits RPC calls per endpoint URL.
type EndpointLimiters struct {
	rps      float64
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
}

func NewEndpointLimiters(rps float64) *EndpointLimiters {
	if rps <= 0 {
		rps = defaultRpcRPS
	}
	return &EndpointLimiters{
		rps:      rps,
		limiters: make(map[string]*rate.Limiter),
	}
}

func (e *EndpointLimiters) limiterFor(url string) *rate.Limiter {
	e.mu.Lock()
	defer e.mu.Unlock()
	if lim, ok := e.limiters[url]; ok {
		return lim
	}
	lim := rate.NewLimiter(rate.Limit(e.rps), 1)
	e.limiters[url] = lim
	return lim
}

// Wait blocks until a request token is available for the given endpoint URL.
func (e *EndpointLimiters) Wait(ctx context.Context, url string) error {
	return e.limiterFor(url).Wait(ctx)
}
