package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
	"go.opentelemetry.io/otel/trace"

	"github.com/isutare412/web-memo/api/internal/tracing"
)

const (
	vectorSize   = 1024
	teiBatchSize = 4
)

type Client struct {
	teiBaseURL     string
	qdrantClient   *qdrant.Client
	collectionName string
	httpClient     *http.Client
}

func NewClient(cfg Config) (*Client, error) {
	qdrantClient, err := qdrant.NewClient(&qdrant.Config{
		Host: cfg.QdrantHost,
		Port: cfg.QdrantPort,
	})
	if err != nil {
		return nil, fmt.Errorf("creating qdrant client: %w", err)
	}

	return &Client{
		teiBaseURL:     cfg.TEIBaseURL,
		qdrantClient:   qdrantClient,
		collectionName: cfg.QdrantCollectionName,
		httpClient:     &http.Client{},
	}, nil
}

func (c *Client) EnsureCollection(ctx context.Context) error {
	ctx, span := tracing.StartSpan(ctx, "embedding.Client.EnsureCollection",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceQdrant))
	defer span.End()

	exists, err := c.qdrantClient.CollectionExists(ctx, c.collectionName)
	if err != nil {
		return fmt.Errorf("checking collection existence: %w", err)
	}
	if exists {
		return nil
	}

	if err := c.qdrantClient.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: c.collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     vectorSize,
			Distance: qdrant.Distance_Cosine,
		}),
	}); err != nil {
		return fmt.Errorf("creating collection: %w", err)
	}

	for _, field := range []string{"memo_id", "owner_id"} {
		if _, err := c.qdrantClient.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
			CollectionName: c.collectionName,
			FieldName:      field,
			FieldType:      qdrant.FieldType_FieldTypeKeyword.Enum(),
		}); err != nil {
			return fmt.Errorf("creating field index for %s: %w", field, err)
		}
	}

	return nil
}

type embedRequest struct {
	Inputs []string `json:"inputs"`
}

func (c *Client) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	var allEmbeddings [][]float32
	for i := 0; i < len(texts); i += teiBatchSize {
		end := min(i+teiBatchSize, len(texts))
		batch, err := c.embedBatch(ctx, texts[i:end])
		if err != nil {
			return nil, fmt.Errorf("embedding batch %d-%d: %w", i, end, err)
		}
		allEmbeddings = append(allEmbeddings, batch...)
	}
	return allEmbeddings, nil
}

func (c *Client) embedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	ctx, span := tracing.StartSpan(ctx, "embedding.Client.embedBatch",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceTEI))
	defer span.End()

	body, err := json.Marshal(embedRequest{Inputs: texts})
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.teiBaseURL+"/embed", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TEI returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var embeddings [][]float32
	if err := json.NewDecoder(resp.Body).Decode(&embeddings); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return embeddings, nil
}

func (c *Client) DeleteByMemoID(ctx context.Context, memoID uuid.UUID) error {
	ctx, span := tracing.StartSpan(ctx, "embedding.Client.DeleteByMemoID",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceQdrant))
	defer span.End()

	_, err := c.qdrantClient.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: c.collectionName,
		Points: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Filter{
				Filter: &qdrant.Filter{
					Must: []*qdrant.Condition{
						qdrant.NewMatch("memo_id", memoID.String()),
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("deleting points: %w", err)
	}

	return nil
}

func (c *Client) UpsertChunks(ctx context.Context, memoID, ownerID uuid.UUID, embeddings [][]float32) error {
	ctx, span := tracing.StartSpan(ctx, "embedding.Client.UpsertChunks",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceQdrant))
	defer span.End()

	points := make([]*qdrant.PointStruct, len(embeddings))
	for i, emb := range embeddings {
		points[i] = &qdrant.PointStruct{
			Id:      qdrant.NewIDUUID(uuid.NewString()),
			Vectors: qdrant.NewVectors(emb...),
			Payload: qdrant.NewValueMap(map[string]any{
				"memo_id":     memoID.String(),
				"owner_id":    ownerID.String(),
				"chunk_index": int64(i),
			}),
		}
	}

	_, err := c.qdrantClient.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: c.collectionName,
		Points:         points,
	})
	if err != nil {
		return fmt.Errorf("upserting points: %w", err)
	}

	return nil
}

func (c *Client) ExistsByMemoID(ctx context.Context, memoID uuid.UUID) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "embedding.Client.ExistsByMemoID",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceQdrant))
	defer span.End()

	count, err := c.qdrantClient.Count(ctx, &qdrant.CountPoints{
		CollectionName: c.collectionName,
		Filter: &qdrant.Filter{
			Must: []*qdrant.Condition{
				qdrant.NewMatch("memo_id", memoID.String()),
			},
		},
		Exact: qdrant.PtrOf(true),
	})
	if err != nil {
		return false, fmt.Errorf("counting points: %w", err)
	}

	return count > 0, nil
}

func (c *Client) Close() error {
	return c.qdrantClient.Close()
}
