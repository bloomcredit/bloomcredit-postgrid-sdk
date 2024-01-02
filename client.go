// Package postgrid provides a client to interact with the PostGrid REST API.
package postgrid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL = "https://api.postgrid.com/v1"
)

// Client allows for interacting with the postgrid api.
type Client struct {
	httpClient *http.Client
	apiKey     string
}

// NewClient constructs a new client with the given api key.
func NewClient(apiKey string, opts ...Option) *Client {
	options := options{
		httpClient: &http.Client{},
	}

	for _, opt := range opts {
		opt.apply(&options)
	}

	return &Client{
		apiKey:     apiKey,
		httpClient: options.httpClient,
	}
}

// VerifyAddress calls the Verify Address endpoint from the postgrid api.
// https://avdocs.postgrid.com/#1061f2ea-00ee-4977-99da-a54872de28c2
func (c *Client) VerifyAddress(ctx context.Context, req any) (any, error) {
	r, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s", baseURL, "/addver/verifications"), req)
	if err != nil {
		return nil, err
	}

	var resp any
	if err = c.send(r, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// BatchVerifyAddresses calls the Batch Verify Address endpoint from the postgrid api.
// https://avdocs.postgrid.com/#94520412-5072-4f5a-a2e2-49981b66a347
func (c *Client) BatchVerifyAddresses(ctx context.Context, req BatchVerifyAddressesRequest) (BatchVerifyAddressesResponse, error) {
	r, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s", baseURL, "/addver/verifications/batch"), req)
	if err != nil {
		return BatchVerifyAddressesResponse{}, err
	}
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
	// Set default headers
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

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

func (c *Client) newRequest(ctx context.Context, method string, url string, body any) (*http.Request, error) {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		buf = bytes.NewBuffer(b)
	}

	return http.NewRequestWithContext(ctx, method, url, buf)
}
