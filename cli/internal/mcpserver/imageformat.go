package mcpserver

import (
	"bytes"
	"encoding/binary"
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

// imageHeadSize is how many leading bytes detectImageFormat can make use of.
const imageHeadSize = 128

var errUnsupportedImage = errors.New("unsupported image format: file content is not JPEG, PNG, WEBP, AVIF or HEIC")

// detectImageFormat identifies the image format from the file's leading bytes
// (the first 12 bytes decide all formats but HEIF-family files, whose ftyp
// box lists compatible brands that may follow; pass up to imageHeadSize bytes).
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
		case "heic", "heix":
			return formatHEIC, nil
		case "mif1", "msf1":
			// Generic HEIF brands: the compatible brands tell AVIF from HEIC.
			return heifFormatFromCompatibleBrands(head), nil
		}
	}
	return "", errUnsupportedImage
}

// heifFormatFromCompatibleBrands scans the compatible brands of the ftyp box
// at the start of head (bounded by the box size and the bytes available).
// The first recognised brand decides; with no hint it is HEIC.
func heifFormatFromCompatibleBrands(head []byte) gen.ImageFormat {
	end := len(head)
	if size := int(binary.BigEndian.Uint32(head[:4])); size >= 16 && size < end {
		end = size
	}
	for i := 16; i+4 <= end; i += 4 {
		switch string(head[i : i+4]) {
		case "avif", "avis":
			return formatAVIF
		case "heic", "heix", "hevc", "hevx":
			return formatHEIC
		}
	}
	return formatHEIC
}
