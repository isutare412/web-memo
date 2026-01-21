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
	imageerClient       *gateway.ClientWithResponses
	projectID           string
	downscalePresetName string
}

func NewService(cfg Config, imageerClient *gateway.ClientWithResponses) *Service {
	return &Service{
		imageerClient:       imageerClient,
		projectID:           cfg.ProjectID,
		downscalePresetName: cfg.DownscalePresetName,
	}
}

func (s *Service) CreateUploadURL(ctx context.Context, fileName string, format images.Format,
) (*model.UploadURL, error) {
	reqBody := gateway.CreateUploadURLJSONRequestBody{
		FileName: fileName,
		Format:   format,
	}
	if s.downscalePresetName != "" {
		reqBody.PresetNames = []string{s.downscalePresetName}
	}

	resp, err := s.imageerClient.CreateUploadURLWithResponse(ctx, s.projectID, reqBody)
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

func (s *Service) GetImage(ctx context.Context, imageID string, waitUntilProcessed bool,
) (*model.Image, error) {
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

	img := &model.Image{
		ID:    resp.JSON200.ID,
		State: resp.JSON200.State,
		Original: &model.ImageData{
			URL:    resp.JSON200.URL,
			Format: resp.JSON200.Format,
		},
	}

	// Find downscaled variant if preset name is configured
	if s.downscalePresetName != "" {
		for _, v := range resp.JSON200.Variants {
			if v.PresetName == s.downscalePresetName && v.State == images.VariantStateReady {
				img.Downscaled = &model.ImageData{
					URL:    v.URL,
					Format: v.Format,
				}
				break
			}
		}
	}

	return img, nil
}
