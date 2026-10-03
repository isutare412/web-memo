package mcpserver

import (
	"strings"
	"testing"
)

func TestDetectImageFormat(t *testing.T) {
	ftyp := func(brand string) []byte {
		return append([]byte{0, 0, 0, 0x18}, []byte("ftyp"+brand+"\x00\x00\x00\x00")...)
	}
	tests := []struct {
		name    string
		head    []byte
		want    string
		wantErr bool
	}{
		{name: "jpeg", head: []byte{0xFF, 0xD8, 0xFF, 0xE0, 0, 0x10}, want: "JPEG"},
		{name: "png", head: []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1A, '\n'}, want: "PNG"},
		{name: "webp", head: []byte("RIFF\x24\x00\x00\x00WEBPVP8 "), want: "WEBP"},
		{name: "avif", head: ftyp("avif"), want: "AVIF"},
		{name: "avis", head: ftyp("avis"), want: "AVIF"},
		{name: "heic", head: ftyp("heic"), want: "HEIC"},
		{name: "heix", head: ftyp("heix"), want: "HEIC"},
		{name: "mif1", head: ftyp("mif1"), want: "HEIC"},
		{name: "msf1", head: ftyp("msf1"), want: "HEIC"},
		{name: "text file", head: []byte("hello, this is not an image"), wantErr: true},
		{name: "wav riff", head: []byte("RIFF\x24\x00\x00\x00WAVEfmt "), wantErr: true},
		{name: "mp4 brand", head: ftyp("isom"), wantErr: true},
		{name: "empty", head: nil, wantErr: true},
		{name: "truncated jpeg", head: []byte{0xFF, 0xD8}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detectImageFormat(tt.head)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("got %q, want error", got)
				}
				if !strings.Contains(err.Error(), "JPEG") {
					t.Errorf("error %q should list supported formats", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("format = %q, want %q", got, tt.want)
			}
		})
	}
}
