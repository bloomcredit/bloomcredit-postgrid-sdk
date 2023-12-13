package postgrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress_MarshalJSON(t *testing.T) {
	type fields struct {
		Line1           string
		Line2           string
		City            string
		ProvinceOrState string
		PostalOrZip     string
		Country         string
		InputID         string
		String          string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success with string",
			fields: fields{
				String: "101 place st, NY",
			},
			want:    []byte(`"101 place st, NY"`),
			wantErr: assert.NoError,
		},
		{
			name: "success no string",
			fields: fields{
				Line1: "101 place st",
				City:  "NY",
			},
			want:    []byte(`{"line1":"101 place st","city":"NY"}`),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Address{
				Line1:           tt.fields.Line1,
				Line2:           tt.fields.Line2,
				City:            tt.fields.City,
				ProvinceOrState: tt.fields.ProvinceOrState,
				PostalOrZip:     tt.fields.PostalOrZip,
				Country:         tt.fields.Country,
				InputID:         tt.fields.InputID,
				String:          tt.fields.String,
			}
			got, err := a.MarshalJSON()
			if !tt.wantErr(t, err, "MarshalJSON()") {
				return
			}
			assert.Equal(t, tt.want, got, "MarshalJSON want = %v, got %v", string(tt.want), string(got))
		})
	}
}
