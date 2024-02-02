package postgrid

import (
	"net/http"

	"golang.org/x/time/rate"
)

type options struct {
	httpClient  *http.Client
	rateLimiter *rate.Limiter
}

// Option represents optional arguments for constructing a postgrid client.
type Option interface {
	apply(*options)
}

type httpClientOption struct {
	client *http.Client
}

func (h httpClientOption) apply(opts *options) {
	opts.httpClient = h.client
}

// WithHTTPClient configures the postgrid client to use the given http.Client.
func WithHTTPClient(client *http.Client) Option {
	return httpClientOption{client: client}
}

type rateLimiterOption struct {
	limiter *rate.Limiter
}

func (r rateLimiterOption) apply(opts *options) {
	opts.rateLimiter = r.limiter
}

// WithRateLimiter configures the postgrid client to use the given rate limiter.
func WithRateLimiter(limiter *rate.Limiter) Option {
	return rateLimiterOption{limiter: limiter}
}
