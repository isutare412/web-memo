package embedding

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type jobType int

const (
	jobEmbed jobType = iota
	jobDelete
)

type job struct {
	typ    jobType
	embed  model.EmbeddingJob
	memoID uuid.UUID
}

type Worker struct {
	client *Client
	jobs   chan job
	wg     sync.WaitGroup
	errs   chan error
}

func NewWorker(cfg Config, client *Client) *Worker {
	return &Worker{
		client: client,
		jobs:   make(chan job, cfg.JobBufferSize),
		errs:   make(chan error, 1),
	}
}

func (w *Worker) EnsureCollection(ctx context.Context) error {
	if err := w.client.EnsureCollection(ctx); err != nil {
		return fmt.Errorf("ensuring qdrant collection: %w", err)
	}
	return nil
}

func (w *Worker) Run() <-chan error {
	go func() {
		ctx := context.Background()

		for j := range w.jobs {
			switch j.typ {
			case jobEmbed:
				w.processEmbed(ctx, j.embed)
			case jobDelete:
				w.processDelete(ctx, j.memoID)
			}
			w.wg.Done()
		}
	}()

	return w.errs
}

func (w *Worker) Shutdown(ctx context.Context) error {
	close(w.jobs)

	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}

	if err := w.client.Close(); err != nil {
		return fmt.Errorf("closing embedding client: %w", err)
	}

	return nil
}

func (w *Worker) Enqueue(ej model.EmbeddingJob) {
	w.wg.Add(1)
	select {
	case w.jobs <- job{typ: jobEmbed, embed: ej}:
	default:
		w.wg.Done()
		slog.Warn("embedding job buffer full, dropping job", "memoId", ej.MemoID)
	}
}

func (w *Worker) EnqueueDelete(memoID uuid.UUID) {
	w.wg.Add(1)
	select {
	case w.jobs <- job{typ: jobDelete, memoID: memoID}:
	default:
		w.wg.Done()
		slog.Warn("embedding job buffer full, dropping delete job", "memoId", memoID)
	}
}

func (w *Worker) processEmbed(ctx context.Context, ej model.EmbeddingJob) {
	text := prepareText(ej.Title, ej.Content)
	chunks := chunkText(text)

	embeddings, err := w.client.Embed(ctx, chunks)
	if err != nil {
		slog.Error("failed to embed memo", "memoId", ej.MemoID, "error", err)
		return
	}

	if err := w.client.DeleteByMemoID(ctx, ej.MemoID); err != nil {
		slog.Error("failed to delete old embeddings", "memoId", ej.MemoID, "error", err)
		return
	}

	if err := w.client.UpsertChunks(ctx, ej.MemoID, ej.OwnerID, embeddings); err != nil {
		slog.Error("failed to upsert embeddings", "memoId", ej.MemoID, "error", err)
		return
	}

	slog.Info("embedded memo", "memoId", ej.MemoID, "chunks", len(chunks))
}

func (w *Worker) processDelete(ctx context.Context, memoID uuid.UUID) {
	if err := w.client.DeleteByMemoID(ctx, memoID); err != nil {
		slog.Error("failed to delete embeddings", "memoId", memoID, "error", err)
		return
	}

	slog.Info("deleted memo embeddings", "memoId", memoID)
}
