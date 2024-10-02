package port

import (
	"context"
	"io"
)

type ObjectStorage interface {
	UploadObject(ctx context.Context, key string, body io.Reader) error
	DeleteObjects(ctx context.Context, keys []string) error
	ListObjectKeysPrefix(ctx context.Context, prefix string) ([]string, error)
}
