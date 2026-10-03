package mcpserver

import (
	"bytes"
	"errors"

	"github.com/isutare412/web-memo/cli/internal/apiclient/gen"
)

// Image formats accepted by the API (gen.ImageFormat is a plain string).
const (
	formatJPEG gen.ImageFormat = "JPEG"
	formatPNG  gen.ImageFormat = "PNG"
	formatWEBP gen.ImageFormat = "WEBP"
	formatAVIF gen.ImageFormat = "AVIF"
	formatHEIC gen.ImageFormat = "HEIC"
)

var errUnsupportedImage = errors.New("unsupported image format: file content is not JPEG, PNG, WEBP, AVIF or HEIC")

// detectImageFormat identifies the image format from the file's leading bytes
// (the first 12 are enough).
func detectImageFormat(head []byte) (gen.ImageFormat, error) {
	switch {
	case bytes.HasPrefix(head, []byte{0xFF, 0xD8, 0xFF}):
		return formatJPEG, nil
	case bytes.HasPrefix(head, []byte{0x89, 'P', 'N', 'G'}):
		return formatPNG, nil
	case len(head) >= 12 && bytes.Equal(head[:4], []byte("RIFF")) && bytes.Equal(head[8:12], []byte("WEBP")):
		return formatWEBP, nil
	case len(head) >= 12 && bytes.Equal(head[4:8], []byte("ftyp")):
		// ISO base media file: the major brand tells AVIF from HEIC.
		switch string(head[8:12]) {
		case "avif", "avis":
			return formatAVIF, nil
		case "heic", "heix", "mif1", "msf1":
			return formatHEIC, nil
		}
	}
	return "", errUnsupportedImage
}
