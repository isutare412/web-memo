package pkgerr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("error only for test")

func TestKnown_Unwrap(t *testing.T) {
	type fields struct {
		Code   Code
		Simple error
		Origin error
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "unwrap_origin",
			fields: fields{
				Origin: errTest,
			},
			wantErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kerr := Known{
				Code:   tt.fields.Code,
				Simple: tt.fields.Simple,
				Origin: tt.fields.Origin,
			}
			err := fmt.Errorf("unwrapped: %w", kerr)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.ErrorAs(t, err, &tt.wantErr)
		})
	}
}
