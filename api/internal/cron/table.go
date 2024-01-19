package cron

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/isutare412/web-memo/api/internal/core/port"
)

type Table struct {
	scheduler   *cron.Cron
	memoService port.MemoService
	lifeErrs    chan error
}

func NewTable(cfg Config, memoService port.MemoService) *Table {
	scheduler := cron.New(cron.WithLogger(cron.DiscardLogger))
	table := &Table{
		scheduler:   scheduler,
		memoService: memoService,
		lifeErrs:    make(chan error, 1),
	}

	scheduler.Schedule(
		cron.Every(cfg.TagCleanupInterval),
		cron.FuncJob(table.runCleanUpTags),
	)

	return table
}

func (t *Table) Run() <-chan error {
	t.scheduler.Start()
	return t.lifeErrs
}

func (t *Table) Shutdown(ctx context.Context) error {
	stopCtx := t.scheduler.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-stopCtx.Done():
	}
	return nil
}

func (t *Table) runCleanUpTags() {
	defer func() {
		if v := recover(); v != nil {
			slog.Error("cron goroutine panicked", "recover", v, "stack", string(debug.Stack()))
			t.lifeErrs <- fmt.Errorf("panic during tag cleanup: %v", v)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deleteCount, err := t.memoService.DeleteOrphanTags(ctx)
	if err != nil {
		slog.Error("failed to delete orphan tags", "error", err)
		return
	}
	if deleteCount > 0 {
		slog.Info("deleted orphan tags", "count", deleteCount)
	}
}
