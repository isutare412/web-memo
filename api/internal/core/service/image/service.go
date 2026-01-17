package image

import (
	"context"
	"fmt"

	"github.com/isutare412/imageer/pkg/gateway"
	"github.com/isutare412/imageer/pkg/images"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type Service struct {
	imageerClient *gateway.ClientWithResponses
	projectID     string
}

func NewService(cfg Config, imageerClient *gateway.ClientWithResponses) *Service {
	return &Service{
		imageerClient: imageerClient,
		projectID:     cfg.ProjectID,
	}
}

func (s *Service) CreateUploadURL(ctx context.Context, fileName string, format images.Format) (*model.UploadURL, error) {
	resp, err := s.imageerClient.CreateUploadURLWithResponse(ctx, s.projectID, gateway.CreateUploadURLJSONRequestBody{
		FileName: fileName,
		Format:   format,
	})
	if err != nil {
		return nil, fmt.Errorf("calling imageer CreateUploadURL: %w", err)
	}

	if resp.JSON200 == nil {
		return nil, imageerError(resp.StatusCode(), resp.JSONDefault)
	}

	return &model.UploadURL{
		ImageID:   resp.JSON200.ImageID,
		URL:       resp.JSON200.URL,
		Headers:   resp.JSON200.Header,
		ExpiresAt: resp.JSON200.ExpiresAt,
	}, nil
}

func (s *Service) GetImage(ctx context.Context, imageID string, waitUntilProcessed bool) (*model.Image, error) {
	params := &gateway.GetImageParams{
		WaitUntilProcessed: lo.EmptyableToPtr(waitUntilProcessed),
	}

	resp, err := s.imageerClient.GetImageWithResponse(ctx, s.projectID, imageID, params)
	if err != nil {
		return nil, fmt.Errorf("calling imageer GetImage: %w", err)
	}

	if resp.JSON200 == nil {
		return nil, imageerError(resp.StatusCode(), resp.JSONDefault)
	}

	return &model.Image{
		ID:     resp.JSON200.ID,
		State:  resp.JSON200.State,
		URL:    resp.JSON200.URL,
		Format: resp.JSON200.Format,
		Variants: lo.Map(resp.JSON200.Variants, func(v gateway.ImageVariant, _ int) model.ImageVariant {
			return model.ImageVariant{
				ID:         v.ID,
				PresetName: v.PresetName,
				URL:        v.URL,
				Format:     v.Format,
				State:      string(v.State),
			}
		}),
	}, nil
}
