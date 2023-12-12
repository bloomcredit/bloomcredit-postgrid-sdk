package postgrid

import "net/http"

type options struct {
	httpClient *http.Client
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
