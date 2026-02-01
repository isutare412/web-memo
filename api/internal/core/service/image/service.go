package image

import (
	"context"
	"fmt"
	"net/http"

	"github.com/isutare412/imageer/pkg/gateway"
	"github.com/isutare412/imageer/pkg/images"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/trace"
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

	resp, err := s.createUploadURL(ctx, s.projectID, reqBody)
	if err != nil {
		return nil, err
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

	resp, err := s.getImage(ctx, imageID, params)
	if err != nil {
		return nil, err
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

func (s *Service) createUploadURL(ctx context.Context, projectID string,
	req gateway.CreateUploadURLJSONRequestBody,
) (*gateway.CreateUploadURLResponse, error) {
	ctx, span := trace.StartSpan(ctx, "image.Service.createUploadURL")
	defer span.End()

	resp, err := s.imageerClient.CreateUploadURLWithResponse(ctx, s.projectID, req, injectTraceContext)
	if err != nil {
		return nil, fmt.Errorf("calling imageer CreateUploadURL: %w", err)
	}
	return resp, nil
}

func (s *Service) getImage(ctx context.Context, imageID string, params *gateway.GetImageParams,
) (*gateway.GetImageResponse, error) {
	ctx, span := trace.StartSpan(ctx, "image.Service.getImage")
	defer span.End()

	resp, err := s.imageerClient.GetImageWithResponse(ctx, s.projectID, imageID, params, injectTraceContext)
	if err != nil {
		return nil, fmt.Errorf("calling imageer GetImage: %w", err)
	}
	return resp, nil
}

func injectTraceContext(ctx context.Context, req *http.Request) error {
	trace.InjectToHeader(ctx, req.Header)
	return nil
}
