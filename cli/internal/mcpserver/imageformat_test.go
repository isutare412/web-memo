package mcpserver

import (
	"strings"
	"testing"
)

func TestDetectImageFormat(t *testing.T) {
	ftyp := func(brand string) []byte {
		return append([]byte{0, 0, 0, 0x18}, []byte("ftyp"+brand+"\x00\x00\x00\x00")...)
	}
	// ftypBrands builds an ftyp box with a major brand and compatible brands.
	ftypBrands := func(major string, compat ...string) []byte {
		b := []byte("\x00\x00\x00\x00ftyp" + major + "\x00\x00\x00\x00")
		for _, c := range compat {
			b = append(b, c...)
		}
		b[3] = byte(len(b))
		return b
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
		{name: "mif1 with avif compatible", head: ftypBrands("mif1", "miaf", "avif"), want: "AVIF"},
		{name: "msf1 with avis compatible", head: ftypBrands("msf1", "avis"), want: "AVIF"},
		{name: "mif1 with heic compatible", head: ftypBrands("mif1", "heic", "miaf"), want: "HEIC"},
		{name: "mif1 with hevc compatible", head: ftypBrands("mif1", "hevc"), want: "HEIC"},
		{name: "mif1 with unknown compatible", head: ftypBrands("mif1", "miaf"), want: "HEIC"},
		{name: "mif1 avif brand beyond box size ignored", head: append(ftypBrands("mif1", "miaf"), "avif"...), want: "HEIC"},
		{name: "mif1 huge box size is capped at available bytes", head: append([]byte{0xFF, 0xFF, 0xFF, 0xFF}, []byte("ftypmif1\x00\x00\x00\x00avif")...), want: "AVIF"},
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
