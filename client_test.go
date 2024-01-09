package postgrid

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_VerifyAddress(t *testing.T) {
	type args struct {
		ctx context.Context
		req VerifyAddressRequest
	}
	tests := []struct {
		name         string
		args         args
		expectations expectations
		resp         response
		want         VerifiedAddress
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "Success - structured address",
			args: args{
				ctx: context.Background(),
				req: VerifyAddressRequest{
					Address: Address{
						Line1:           "251 e 13th st frnt a",
						Line2:           "",
						City:            "New York",
						ProvinceOrState: "NY",
						PostalOrZip:     "10003",
						Country:         "US",
					},
				},
			},
			expectations: expectations{
				Path:   "/addver/verifications?geocode=true&includeDetails=true",
				Method: http.MethodPost,
				Headers: map[string][]string{
					"X-API-Key": nil,
				},
				Body: VerifyAddressRequest{
					Address: Address{
						Line1:           "251 e 13th st frnt a",
						Line2:           "",
						City:            "New York",
						ProvinceOrState: "NY",
						PostalOrZip:     "10003",
						Country:         "US",
					},
				}.Encode(),
			},
			resp: response{
				Body: Response{
					Status:  ResponseStatusSuccess,
					Message: "Address verification processed",
					Data: mustMarshalJSON(t, VerifiedAddress{
						Line1:           "251 E 13TH ST FRNT A",
						Line2:           "",
						City:            "NEW YORK",
						ProvinceOrState: "NY",
						PostalOrZip:     "10003",
						ZipPlus4:        "5646",
						FirmName:        "",
						Country:         "us",
						Errors:          map[string][]string{},
						Status:          "corrected",
						Details: VerifiedAddressDetails{
							StreetName: "13TH",
							StreetType: "ST",
							County:     "NEW YORK",
						},
						GeocodeResult: GeocodeResult{
							Location: GeocodeLocation{
								Latitude:  40.731862,
								Longitude: -73.985679,
							},
							Accuracy:     1,
							AccuracyType: "rooftop",
						},
					}),
				},
				Status: http.StatusOK,
			},
			want: VerifiedAddress{
				Line1:           "251 E 13TH ST FRNT A",
				Line2:           "",
				City:            "NEW YORK",
				ProvinceOrState: "NY",
				PostalOrZip:     "10003",
				ZipPlus4:        "5646",
				FirmName:        "",
				Country:         "us",
				Errors:          map[string][]string{},
				Status:          "corrected",
				Details: VerifiedAddressDetails{
					StreetName: "13TH",
					StreetType: "ST",
					County:     "NEW YORK",
				},
				GeocodeResult: GeocodeResult{
					Location: GeocodeLocation{
						Latitude:  40.731862,
						Longitude: -73.985679,
					},
					Accuracy:     1,
					AccuracyType: "rooftop",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Success - freeform address",
			args: args{
				ctx: context.Background(),
				req: VerifyAddressRequest{
					Address: Address{
						String: "251 e 13th st frnt a, New York, NY 10003",
					},
				},
			},
			expectations: expectations{
				Path:   "/addver/verifications?geocode=true&includeDetails=true",
				Method: http.MethodPost,
				Headers: map[string][]string{
					"X-API-Key": nil,
				},
				Body: VerifyAddressRequest{
					Address: Address{
						String: "251 e 13th st frnt a, New York, NY 10003",
					},
				}.Encode(),
			},
			resp: response{
				Body: Response{
					Status:  ResponseStatusSuccess,
					Message: "Address verification processed",
					Data: mustMarshalJSON(t, VerifiedAddress{
						Line1:           "251 E 13TH ST FRNT A",
						Line2:           "",
						City:            "NEW YORK",
						ProvinceOrState: "NY",
						PostalOrZip:     "10003",
						ZipPlus4:        "5646",
						FirmName:        "",
						Country:         "us",
						Errors:          map[string][]string{},
						Status:          "corrected",
						Details: VerifiedAddressDetails{
							StreetName: "13TH",
							StreetType: "ST",
							County:     "NEW YORK",
						},
						GeocodeResult: GeocodeResult{
							Location: GeocodeLocation{
								Latitude:  40.731862,
								Longitude: -73.985679,
							},
							Accuracy:     1,
							AccuracyType: "rooftop",
						},
					}),
				},
				Status: http.StatusOK,
			},
			want: VerifiedAddress{
				Line1:           "251 E 13TH ST FRNT A",
				Line2:           "",
				City:            "NEW YORK",
				ProvinceOrState: "NY",
				PostalOrZip:     "10003",
				ZipPlus4:        "5646",
				FirmName:        "",
				Country:         "us",
				Errors:          map[string][]string{},
				Status:          "corrected",
				Details: VerifiedAddressDetails{
					StreetName: "13TH",
					StreetType: "ST",
					County:     "NEW YORK",
				},
				GeocodeResult: GeocodeResult{
					Location: GeocodeLocation{
						Latitude:  40.731862,
						Longitude: -73.985679,
					},
					Accuracy:     1,
					AccuracyType: "rooftop",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "server error",
			args: args{
				ctx: context.Background(),
				req: VerifyAddressRequest{
					Address: Address{
						String: "251 e 13th st frnt a, New York, NY 10003",
					},
				},
			},
			expectations: expectations{
				Path:   "/addver/verifications?geocode=true&includeDetails=true",
				Method: http.MethodPost,
				Headers: map[string][]string{
					"X-API-Key": nil,
				},
				Body: VerifyAddressRequest{
					Address: Address{
						String: "251 e 13th st frnt a, New York, NY 10003",
					},
				}.Encode(),
			},
			resp: response{
				Body: Response{
					Status:  ResponseStatusError,
					Message: "Error processing address verification",
					Data:    nil,
				},
				Status: http.StatusInternalServerError,
			},
			want:    VerifiedAddress{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := startTestServer(t, tt.expectations, tt.resp)
			t.Cleanup(srv.Close)
			client := NewClient("", srv.URL, WithHTTPClient(srv.Client()))

			got, err := client.VerifyAddress(tt.args.ctx, tt.args.req)
			if !tt.wantErr(t, err, "VerifyAddress(%v, %v)", tt.args.ctx, tt.args.req) {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

type expectations struct {
	Path    string
	Method  string
	Headers http.Header
	Body    io.Reader
}

type response struct {
	Body   any
	Status int
}

func startTestServer(tb testing.TB, expectations expectations, resp response) *httptest.Server {
	tb.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifyExpectations(tb, r, expectations)

		buf, err := json.Marshal(resp.Body)
		require.NoError(tb, err)
		w.WriteHeader(resp.Status)
		_, err = w.Write(buf)
		require.NoError(tb, err)
	}))

	return srv
}

func verifyExpectations(tb testing.TB, r *http.Request, expectations expectations) {
	assert.Equal(tb, expectations.Path, r.URL.String())
	assert.Equal(tb, expectations.Method, r.Method)
	for k, v := range expectations.Headers {
		assert.Equal(tb, v, r.Header[k])
	}
	buf, err := io.ReadAll(r.Body)
	require.NoError(tb, err)
	expectedBody, err := io.ReadAll(expectations.Body)
	require.NoError(tb, err)
	assert.Equal(tb, expectedBody, buf)
}

func mustMarshalJSON(tb testing.TB, body any) json.RawMessage {
	tb.Helper()

	b, err := json.Marshal(body)
	require.NoError(tb, err)

	return b
}

// {
// 	"status": "success",
// 	"message": "Address verification processed.",
// 	"data": {
// 	  "city": "NEW YORK",
// 	  "country": "us",
// 	  "countryName": "UNITED STATES",
// 	  "details": {
// 		"streetName": "13TH",
// 		"streetType": "ST",
// 		"streetDirection": "E",
// 		"preDirection": "E",
// 		"streetNumber": "251",
// 		"suiteID": "A",
// 		"suiteKey": "FRNT",
// 		"county": "NEW YORK"
// 	  },
// 	  "errors": {},
// 	  "firmName": "MILK BAR",
// 	  "geocodeResult": {
// 		"location": {
// 		  "lat": 40.731862,
// 		  "lng": -73.985679
// 		},
// 		"accuracy": 1,
// 		"accuracyType": "rooftop"
// 	  },
// 	  "line1": "251 E 13TH ST FRNT A",
// 	  "postalOrZip": "10003",
// 	  "provinceOrState": "NY",
// 	  "status": "corrected",
// 	  "zipPlus4": "5646"
// 	}
//   }
