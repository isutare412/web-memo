package repeatjob

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"

	"github.com/isutare412/web-memo/api/internal/core/port"
)

type Trigger struct {
	memoService port.MemoService

	tagCleanupInterval time.Duration

	waitGroup     *sync.WaitGroup
	lifeErrs      chan error
	lifeCtx       context.Context
	lifeCtxCancel context.CancelFunc
}

func NewTrigger(cfg Config, memoService port.MemoService) *Trigger {
	lifeCtx, lifeCtxCancel := context.WithCancel(context.Background())

	return &Trigger{
		memoService:        memoService,
		tagCleanupInterval: cfg.TagCleanupInterval,
		waitGroup:          &sync.WaitGroup{},
		lifeErrs:           make(chan error, 1),
		lifeCtx:            lifeCtx,
		lifeCtxCancel:      lifeCtxCancel,
	}
}

func (t *Trigger) Run() <-chan error {
	t.runCleanUpTags()
	return t.lifeErrs
}

func (t *Trigger) Shutdown(ctx context.Context) error {
	t.lifeCtxCancel()

	workerDone := make(chan struct{})
	go func() {
		t.waitGroup.Wait()
		close(workerDone)
	}()

	select {
	case <-workerDone:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("timed out waiting for workers to close")
	}
}

func (t *Trigger) runCleanUpTags() {
	t.waitGroup.Add(1)
	go func() {
		defer t.waitGroup.Done()
		defer func() {
			if v := recover(); v != nil {
				slog.Error("goroutine panicked", "recover", v, "stack", string(debug.Stack()))
				t.lifeErrs <- fmt.Errorf("panic during tag cleanup: %v", v)
			}
		}()

		ticker := time.NewTicker(t.tagCleanupInterval)
		defer ticker.Stop()

		for {
			if err := t.cleanUpTags(context.Background()); err != nil {
				slog.Error("failed to clean up tags", "error", err)
			}

			select {
			case <-t.lifeCtx.Done():
				slog.Info("stop tag cleanup job")
				return
			case <-ticker.C:
			}
		}
	}()
}

func (t *Trigger) cleanUpTags(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	deleteCount, err := t.memoService.DeleteOrphanTags(ctx)
	if err != nil {
		return fmt.Errorf("deleting orhpan tags: %w", err)
	}
	if deleteCount > 0 {
		slog.Info("deleted orphan tags", "count", deleteCount)
	}

	return nil
}
