package postgrid

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"
	"strings"
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

// VerifyAddressRequest represents the request model to be sent to the Verify Address endpoint.
type VerifyAddressRequest struct {
	Address Address
}

func (r VerifyAddressRequest) Encode() io.Reader {
	data := url.Values{}
	if r.Address.String != "" {
		data.Add("address", r.Address.String)
	} else {
		data.Add("address[line1]", r.Address.Line1)
		data.Add("address[line2]", r.Address.Line2)
		data.Add("address[city]", r.Address.City)
		data.Add("address[provinceOrState]", r.Address.ProvinceOrState)
		data.Add("address[postalOrZip]", r.Address.PostalOrZip)
		data.Add("address[country]", r.Address.Country)
	}

	return strings.NewReader(data.Encode())
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
	Location     GeocodeLocation `json:"location"`
	Accuracy     float32         `json:"accuracy"`
	AccuracyType string          `json:"accuracyType"`
}

type GeocodeLocation struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

// VerifiedAddress represents an address that has been verified by postgrid.
type VerifiedAddress struct {
	Line1           string                 `json:"line1"`
	Line2           string                 `json:"line2"`
	City            string                 `json:"city"`
	ProvinceOrState string                 `json:"provinceOrState"`
	PostalOrZip     string                 `json:"postalOrZip"`
	ZipPlus4        string                 `json:"zipPlus4"`
	FirmName        string                 `json:"firmName"`
	Country         string                 `json:"country"`
	Errors          map[string][]string    `json:"errors"`
	Status          string                 `json:"status"`
	Details         VerifiedAddressDetails `json:"details"`
	GeocodeResult   GeocodeResult          `json:"geocodeResult"`
}
type VerifiedAddressDetails struct {
	StreetName                         string `json:"streetName"`
	StreetType                         string `json:"streetType"`
	StreetDirection                    string `json:"streetDirection"`
	StreetNumber                       string `json:"streetNumber"`
	SuiteID                            string `json:"suiteID"`
	SuiteKey                           string `json:"suiteKey"`
	BoxID                              string `json:"boxID"`
	DeliveryInstallationAreaName       string `json:"deliveryInstallationAreaName"`
	DeliveryInstallationType           string `json:"deliveryInstallationType"`
	DeliveryInstallationQualifier      string `json:"deliveryInstallationQualifier"`
	RuralRouteNumber                   string `json:"ruralRouteNumber"`
	RuralRouteType                     string `json:"ruralRouteType"`
	ExtraInfo                          string `json:"extraInfo"`
	County                             string `json:"county"`
	CountyNumber                       string `json:"countyNum"`
	Residential                        bool   `json:"residential"`
	Vacant                             bool   `json:"vacant"`
	PreDirection                       string `json:"preDirection"`
	PostDirection                      string `json:"postDirection"`
	USCensusCMSA                       string `json:"usCensusCMSA"`
	USCensusBlockNumber                string `json:"usCensusBlockNumber"`
	USCensusTractNumber                string `json:"usCensusTractNumber"`
	USCongressionalDistrictNumber      string `json:"usCongressionalDistrictNumber"`
	USCensusMA                         string `json:"usCensusMA"`
	USCensusMSA                        string `json:"usCensusMSA"`
	USCensusPMSA                       string `json:"usCensusPMSA"`
	USHasDaylightSavings               bool   `json:"usHasDaylightSavings"`
	USIntelligentMailBarcodeKey        string `json:"usIntelligentMailBarcodeKey"`
	USMailingsCSKey                    string `json:"usMailingsCSKey"`
	USMailingsCarrierRoute             string `json:"usMailingsCarrierRoute"`
	USMailingsCheckDigit               string `json:"usMailingsCheckDigit"`
	USMailingsDefaultFlag              bool   `json:"usMailingsDefaultFlag"`
	USMailingsDeliveryPoint            string `json:"usMailingsDeliveryPoint"`
	USMailingsDpvConfirmationIndicator string `json:"usMailingsDpvConfirmationIndicator"`
	USMailingsDpvCrmaIndicator         string `json:"usMailingsDpvCrmaIndicator"`
	USMailingsDpvFootnote1             string `json:"usMailingsDpvFootnote1"`
	USMailingsDpvFootnote2             string `json:"usMailingsDpvFootnote2"`
	USMailingsDpvFootnote3             string `json:"usMailingsDpvFootnote3"`
	USMailingsElotAscDesc              string `json:"usMailingsElotAscDesc"`
	USMailingsElotSequenceNumber       string `json:"usMailingsElotSequenceNumber"`
	USMailingsEWSFlag                  bool   `json:"usMailingsEWSFlag"`
	USMailingsLACSFlag                 string `json:"usMailingsLACSFlag"`
	USMailingsLACSReturnCode           string `json:"usMailingsLACSReturnCode"`
}

// VerifiedAddressResponse ...
type VerifiedAddressResponse struct {
	VerifiedAddress VerifiedAddress `json:"verifiedAddress"`
}
