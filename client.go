// Package postgrid provides a client to interact with the PostGrid REST API.
package postgrid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

const (
	BaseURL = "https://api.postgrid.com/v1"
)

// Client allows for interacting with the postgrid api.
type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string

	rateLimiter *rate.Limiter
}

// NewClient constructs a new client with the given api key.
func NewClient(apiKey string, baseURL string, opts ...Option) *Client {
	options := options{
		httpClient:  &http.Client{},
		rateLimiter: rate.NewLimiter(5, 5),
	}

	for _, opt := range opts {
		opt.apply(&options)
	}

	return &Client{
		apiKey:      apiKey,
		baseURL:     baseURL,
		httpClient:  options.httpClient,
		rateLimiter: options.rateLimiter,
	}
}

// VerifyAddress calls the Verify Address endpoint from the postgrid api.
// https://avdocs.postgrid.com/#1061f2ea-00ee-4977-99da-a54872de28c2
func (c *Client) VerifyAddress(ctx context.Context, req VerifyAddressRequest) (VerifiedAddress, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.baseURL, "/addver/verifications"), req.Encode())
	if err != nil {
		return VerifiedAddress{}, err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	params := r.URL.Query()
	params.Set("includeDetails", "true")
	params.Set("geocode", "true")
	r.URL.RawQuery = params.Encode()

	var resp VerifiedAddress
	if err = c.send(r, &resp); err != nil {
		return VerifiedAddress{}, err
	}

	return resp, nil
}

// BatchVerifyAddresses calls the Batch Verify Address endpoint from the postgrid api.
// https://avdocs.postgrid.com/#94520412-5072-4f5a-a2e2-49981b66a347
func (c *Client) BatchVerifyAddresses(ctx context.Context, req BatchVerifyAddressesRequest) (BatchVerifyAddressesResponse, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return BatchVerifyAddressesResponse{}, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.baseURL, "/addver/verifications/batch"), bytes.NewBuffer(reqJSON))
	if err != nil {
		return BatchVerifyAddressesResponse{}, err
	}
	r.Header.Set("Content-Type", "application/json")
	params := r.URL.Query()
	params.Set("includeDetails", "true")
	params.Set("geocode", "true")
	r.URL.RawQuery = params.Encode()

	var resp BatchVerifyAddressesResponse
	if err = c.send(r, &resp); err != nil {
		return BatchVerifyAddressesResponse{}, err
	}

	return resp, nil
}

// send initiates the http request and unmarshals the response into the object passed in.
func (c *Client) send(req *http.Request, v any) error {
	// Respect rate limit
	if err := c.rateLimiter.Wait(req.Context()); err != nil {
		return err
	}

	// Set default headers
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("decoding response envelope: %w", err)
	}

	if response.Status == ResponseStatusError {
		return fmt.Errorf("postgrid error: %s", response.Message)
	}

	if v == nil {
		return nil
	}

	return json.Unmarshal(response.Data, v)
}
