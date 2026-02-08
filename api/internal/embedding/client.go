package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"

	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/tracing"
)

const (
	vectorSize   = 1024
	teiBatchSize = 4
	rrfK         = 2
)

type Client struct {
	teiBaseURL     string
	bm25BaseURL    string
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
		bm25BaseURL:    cfg.BM25BaseURL,
		qdrantClient:   qdrantClient,
		collectionName: cfg.QdrantCollectionName,
		httpClient:     &http.Client{},
	}, nil
}

func (c *Client) Close() error {
	return c.qdrantClient.Close()
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
		VectorsConfig: qdrant.NewVectorsConfigMap(map[string]*qdrant.VectorParams{
			"dense": {Size: vectorSize, Distance: qdrant.Distance_Cosine},
		}),
		SparseVectorsConfig: qdrant.NewSparseVectorsConfig(map[string]*qdrant.SparseVectorParams{
			"sparse": {Modifier: qdrant.Modifier_Idf.Enum()},
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

func (c *Client) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	ctx, span := tracing.StartSpan(ctx, "embedding.Client.Embed",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceTEI))
	defer span.End()

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

func (c *Client) EmbedSparse(ctx context.Context, texts []string) ([]sparseVector, error) {
	ctx, span := tracing.StartSpan(ctx, "embedding.Client.EmbedSparse",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceBM25))
	defer span.End()

	var allEmbeddings []sparseVector
	for i := 0; i < len(texts); i += teiBatchSize {
		end := min(i+teiBatchSize, len(texts))
		batch, err := c.embedSparseBatch(ctx, texts[i:end])
		if err != nil {
			return nil, fmt.Errorf("sparse embedding batch %d-%d: %w", i, end, err)
		}
		allEmbeddings = append(allEmbeddings, batch...)
	}
	return allEmbeddings, nil
}

func (c *Client) embedSparseBatch(ctx context.Context, texts []string) ([]sparseVector, error) {
	ctx, span := tracing.StartSpan(ctx, "embedding.Client.embedSparseBatch",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceBM25))
	defer span.End()

	body, err := json.Marshal(embedRequest{Inputs: texts})
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.bm25BaseURL+"/embed-sparse", bytes.NewReader(body))
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
		return nil, fmt.Errorf("BM25 returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var responses []sparseEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&responses); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	embeddings := make([]sparseVector, len(responses))
	for i, r := range responses {
		embeddings[i] = sparseVector(r)
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

func (c *Client) UpsertChunks(ctx context.Context, memoID, ownerID uuid.UUID, denseEmbeddings [][]float32, sparseEmbeddings []sparseVector) error {
	ctx, span := tracing.StartSpan(ctx, "embedding.Client.UpsertChunks",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceQdrant))
	defer span.End()

	points := make([]*qdrant.PointStruct, len(denseEmbeddings))
	for i, emb := range denseEmbeddings {
		points[i] = &qdrant.PointStruct{
			Id: qdrant.NewIDUUID(uuid.NewString()),
			Vectors: qdrant.NewVectorsMap(map[string]*qdrant.Vector{
				"dense":  qdrant.NewVectorDense(emb),
				"sparse": qdrant.NewVectorSparse(sparseEmbeddings[i].Indices, sparseEmbeddings[i].Values),
			}),
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

func (c *Client) Search(ctx context.Context, query string, ownerIDFilter *uuid.UUID,
	memoIDFilter []uuid.UUID, scoreThreshold float32, limit int,
) ([]model.SearchResult, error) {
	ctx, span := tracing.StartSpan(ctx, "embedding.Client.Search",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceQdrant))
	defer span.End()

	// Qwen3-Embedding models produce better retrieval vectors when queries
	// include an instruction prefix. Documents are embedded without it.
	// See https://huggingface.co/Qwen/Qwen3-Embedding-0.6B
	instructQuery := "Instruct: Given a search query, retrieve relevant memos that match the query\nQuery:" + query

	eg, eggCtx := errgroup.WithContext(ctx)
	var denseEmbeddings [][]float32
	var sparseEmbeddings []sparseVector

	eg.Go(func() error {
		var err error
		denseEmbeddings, err = c.Embed(eggCtx, []string{instructQuery})
		return err
	})
	eg.Go(func() error {
		var err error
		sparseEmbeddings, err = c.EmbedSparse(eggCtx, []string{query})
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("embedding query: %w", err)
	}

	var conditions []*qdrant.Condition
	if ownerIDFilter != nil {
		conditions = append(conditions, qdrant.NewMatch("owner_id", ownerIDFilter.String()))
	}
	if len(memoIDFilter) > 0 {
		keywords := make([]string, len(memoIDFilter))
		for i, id := range memoIDFilter {
			keywords[i] = id.String()
		}
		conditions = append(conditions, qdrant.NewMatchKeywords("memo_id", keywords...))
	}

	filter := &qdrant.Filter{
		Must: conditions,
	}

	// Run semantic and BM25 searches concurrently.
	var semanticGroups, bm25Groups []*qdrant.PointGroup

	eg2, egCtx2 := errgroup.WithContext(ctx)
	eg2.Go(func() error {
		groups, err := c.qdrantClient.QueryGroups(egCtx2, &qdrant.QueryPointGroups{
			CollectionName: c.collectionName,
			Query:          qdrant.NewQueryDense(denseEmbeddings[0]),
			Using:          qdrant.PtrOf("dense"),
			GroupBy:        "memo_id",
			GroupSize:      qdrant.PtrOf(uint64(1)),
			Limit:          qdrant.PtrOf(uint64(limit)),
			Filter:         filter,
		})
		if err != nil {
			return fmt.Errorf("semantic search: %w", err)
		}
		semanticGroups = groups
		return nil
	})
	eg2.Go(func() error {
		groups, err := c.qdrantClient.QueryGroups(egCtx2, &qdrant.QueryPointGroups{
			CollectionName: c.collectionName,
			Query:          qdrant.NewQuerySparse(sparseEmbeddings[0].Indices, sparseEmbeddings[0].Values),
			Using:          qdrant.PtrOf("sparse"),
			GroupBy:        "memo_id",
			GroupSize:      qdrant.PtrOf(uint64(1)),
			Limit:          qdrant.PtrOf(uint64(limit)),
			ScoreThreshold: qdrant.PtrOf(float32(math.SmallestNonzeroFloat32)),
			Filter:         filter,
		})
		if err != nil {
			return fmt.Errorf("bm25 search: %w", err)
		}
		bm25Groups = groups
		return nil
	})
	if err := eg2.Wait(); err != nil {
		return nil, fmt.Errorf("searching groups: %w", err)
	}

	// Parse semantic results: memoID → (score, rank).
	type rankScore struct {
		score float32
		rank  int
	}
	semanticMap := make(map[uuid.UUID]rankScore, len(semanticGroups))
	for rank, group := range semanticGroups {
		memoIDStr := group.GetId().GetStringValue()
		memoID, err := uuid.Parse(memoIDStr)
		if err != nil {
			return nil, fmt.Errorf("parsing semantic memo_id %q: %w", memoIDStr, err)
		}
		var score float32
		if hits := group.GetHits(); len(hits) > 0 {
			score = hits[0].GetScore()
		}
		semanticMap[memoID] = rankScore{score: score, rank: rank + 1}
	}

	// Parse BM25 results: memoID → (score, rank).
	bm25Map := make(map[uuid.UUID]rankScore, len(bm25Groups))
	for rank, group := range bm25Groups {
		memoIDStr := group.GetId().GetStringValue()
		memoID, err := uuid.Parse(memoIDStr)
		if err != nil {
			return nil, fmt.Errorf("parsing bm25 memo_id %q: %w", memoIDStr, err)
		}
		var score float32
		if hits := group.GetHits(); len(hits) > 0 {
			score = hits[0].GetScore()
		}
		bm25Map[memoID] = rankScore{score: score, rank: rank + 1}
	}

	// Collect all unique memo IDs.
	allIDs := make(map[uuid.UUID]struct{}, len(semanticMap)+len(bm25Map))
	for id := range semanticMap {
		allIDs[id] = struct{}{}
	}
	for id := range bm25Map {
		allIDs[id] = struct{}{}
	}

	// Compute RRF scores.
	var results []model.SearchResult
	for memoID := range allIDs {
		var rrfScore float32
		var semanticScore, bm25Score float32

		if d, ok := semanticMap[memoID]; ok {
			rrfScore += 1.0 / float32(rrfK+d.rank)
			semanticScore = d.score
		}
		if s, ok := bm25Map[memoID]; ok {
			rrfScore += 1.0 / float32(rrfK+s.rank)
			bm25Score = s.score
		}

		if rrfScore < scoreThreshold {
			continue
		}

		results = append(results, model.SearchResult{
			MemoID:        memoID,
			RRFScore:      rrfScore,
			SemanticScore: semanticScore,
			BM25Score:     bm25Score,
		})
	}

	// Sort by RRF > Semantic > BM25 score descending.
	sort.Slice(results, func(i, j int) bool {
		a, b := results[i], results[j]
		if a.RRFScore != b.RRFScore {
			return a.RRFScore > b.RRFScore
		}
		if a.SemanticScore != b.SemanticScore {
			return a.SemanticScore > b.SemanticScore
		}
		return a.BM25Score > b.BM25Score
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}
