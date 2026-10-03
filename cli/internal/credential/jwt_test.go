package credential

import (
	"encoding/base64"
	"testing"
	"time"
)

func makeJWT(payload string) string {
	enc := base64.RawURLEncoding
	return enc.EncodeToString([]byte(`{"alg":"none"}`)) + "." + enc.EncodeToString([]byte(payload)) + ".sig"
}

func TestTokenExpiry(t *testing.T) {
	tests := []struct {
		name   string
		token  string
		want   time.Time
		wantOK bool
	}{
		{"valid", makeJWT(`{"exp":1893456000}`), time.Unix(1893456000, 0), true},
		{"garbage", "garbage", time.Time{}, false},
		{"three parts but not base64", "a.!!!.c", time.Time{}, false},
		{"payload not json", makeJWT(`not json`), time.Time{}, false},
		{"missing exp", makeJWT(`{"sub":"x"}`), time.Time{}, false},
		{"empty", "", time.Time{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := TokenExpiry(tt.token)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !got.Equal(tt.want) {
				t.Fatalf("expiry = %v, want %v", got, tt.want)
			}
		})
	}
}
