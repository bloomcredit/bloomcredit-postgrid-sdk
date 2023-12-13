package postgrid

import (
	"encoding/json"
	"strconv"
)

// All possible values for ResponseStatus.
const (
	ResponseStatusSuccess = "success"
	ResponseStatusError   = "error"
)

// MaxBatchSize is the max size for a batch address verification request.
const MaxBatchSize = 2000

// Response represents the generic response wrapper from the postgrid api.
type Response struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// Address represents an address that should be sent for verification.
type Address struct {
	Line1           string `json:"line1,omitempty"`
	Line2           string `json:"line2,omitempty"`
	City            string `json:"city,omitempty"`
	ProvinceOrState string `json:"provinceOrState,omitempty"`
	PostalOrZip     string `json:"postalOrZip,omitempty"`
	Country         string `json:"country,omitempty"`
	InputID         string `json:"inputID,omitempty"`

	// String is used if sending the string representation of the address.
	String string `json:"-"`
}

func (a Address) MarshalJSON() ([]byte, error) {
	if a.String != "" {
		return []byte(strconv.Quote(a.String)), nil
	}

	type Alias Address
	return json.Marshal(Alias(a))
}

// BatchVerifyAddressesRequest represents the request model to be sent to the Batch Veryify Addresses endpoint.
type BatchVerifyAddressesRequest struct {
	Addresses []Address `json:"addresses"`
}

// BatchVerifyAddressesResponse represents the response model from the Batch Verify Addresses endpoint.
type BatchVerifyAddressesResponse struct {
	Results []VerifiedAddressResponse `json:"results"`
}

// GeocodeResult represents ...
type GeocodeResult struct {
	Location     map[string]any
	Accuracy     float32 `json:"accuracy"`
	AccuracyType string  `json:"accuracyType"`
}

// VerifiedAddress represents an address that has been verified by postgrid.
type VerifiedAddress struct {
	Line1           string         `json:"line1"`
	Line2           string         `json:"line2"`
	City            string         `json:"city"`
	ProvinceOrState string         `json:"provinceOrState"`
	PostalOrZip     string         `json:"postalOrZip"`
	ZipPlus4        string         `json:"zipPlus4"`
	FirmName        string         `json:"firmName"`
	Country         string         `json:"country"`
	Errors          any            `json:"errors"`
	Status          string         `json:"status"`
	Details         map[string]any `json:"details"`
	GeocodeResult   GeocodeResult  `json:"geocodeResult"`
}

// VerifiedAddressResponse ...
type VerifiedAddressResponse struct {
	VerifiedAddress VerifiedAddress `json:"verifiedAddress"`
}
