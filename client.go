// Package postgrid provides a client to interact with the PostGrid REST API.
package postgrid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
func (c *Client) VerifyAddress(ctx context.Context, req VerifyAddressRequest) (VerifiedAddress, error) {
	data := url.Values{}
	if req.Address.String != "" {
		data.Add("address", req.Address.String)
	} else {
		data.Add("address[line1]", req.Address.Line1)
		data.Add("address[line2]", req.Address.Line2)
		data.Add("address[city]", req.Address.City)
		data.Add("address[provinceOrState]", req.Address.ProvinceOrState)
		data.Add("address[postalOrZip]", req.Address.PostalOrZip)
		data.Add("address[country]", req.Address.Country)
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", baseURL, "/addver/verifications"), strings.NewReader(data.Encode()))
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

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", baseURL, "/addver/verifications/batch"), bytes.NewBuffer(reqJSON))
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

	fmt.Printf("%s\n", string(response.Data))

	if v == nil {
		return nil
	}

	return json.Unmarshal(response.Data, v)
}
