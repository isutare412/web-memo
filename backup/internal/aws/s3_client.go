package aws

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client struct {
	client *s3.Client
	bucket string
}

func NewS3Client(cfg S3Config) (*S3Client, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.Secret, ""),
		),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("loading aws config: %w", err)
	}

	return &S3Client{
		client: s3.NewFromConfig(awsCfg),
		bucket: cfg.Bucket,
	}, nil
}

func (c *S3Client) UploadObject(ctx context.Context, key string, body io.Reader) error {
	if _, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
		Body:   body,
	}); err != nil {
		return fmt.Errorf("putting object: %w", err)
	}
	return nil
}

func (c *S3Client) DeleteObjects(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	var (
		chunkedKeys = chunkSlice(keys, 1000) // S3 supports bulk delete up to 1000 keys.
		totalErr    error
	)
	for _, keys := range chunkedKeys {
		deleteTargets := make([]types.ObjectIdentifier, 0, len(keys))
		for _, k := range keys {
			deleteTargets = append(deleteTargets, types.ObjectIdentifier{
				Key: &k,
			})
		}

		if _, err := c.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: &c.bucket,
			Delete: &types.Delete{
				Objects: deleteTargets,
			},
		}); err != nil {
			totalErr = errors.Join(totalErr, fmt.Errorf("deleting objects: %w", err))
			continue
		}
	}

	return totalErr
}

func (c *S3Client) ListObjectKeysPrefix(ctx context.Context, prefix string) ([]string, error) {
	var (
		keys              []string
		continuationToken *string
	)
	for {
		output, err := c.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:            &c.bucket,
			Delimiter:         aws.String("/"),
			Prefix:            &prefix,
			ContinuationToken: continuationToken,
		})
		if err != nil {
			return nil, fmt.Errorf("listing objects: %w", err)
		}

		for _, c := range output.Contents {
			keys = append(keys, *c.Key)
		}

		if output.IsTruncated != nil && *output.IsTruncated {
			continuationToken = output.NextContinuationToken
			continue
		}

		return keys, nil
	}
}
