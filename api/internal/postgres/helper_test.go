package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pointerIfNotZero(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		wantNil bool
	}{
		{
			name:    "zero_string",
			value:   "",
			wantNil: true,
		},
		{
			name:    "non_zero_string",
			value:   "foo",
			wantNil: false,
		},
		{
			name:    "zero_int",
			value:   0,
			wantNil: true,
		},
		{
			name:    "non_zero_int",
			value:   42,
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := pointerIfNotZero(tt.value)
			assert.True(t, (ptr == nil) == tt.wantNil)
		})
	}
}
