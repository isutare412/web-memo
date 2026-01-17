package imageer

import (
	"context"
	"net/http"

	"github.com/isutare412/imageer/pkg/gateway"
)

const apiKeyHeader = "X-API-KEY"

func NewClient(cfg Config) (*gateway.ClientWithResponses, error) {
	return gateway.NewClientWithResponses(
		cfg.BaseURL,
		gateway.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			req.Header.Set(apiKeyHeader, cfg.APIKey)
			return nil
		}),
	)
}
