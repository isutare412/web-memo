package cron

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/go-co-op/gocron/v2"

	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/tracing"
)

type Scheduler struct {
	client      gocron.Scheduler
	memoService port.MemoService
	lifeErrs    chan error
}

func NewScheduler(cfg Config, memoService port.MemoService) (*Scheduler, error) {
	client, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	scheduler := &Scheduler{
		client:      client,
		memoService: memoService,
		lifeErrs:    make(chan error, 1),
	}

	client.NewJob(
		gocron.DurationJob(cfg.TagCleanupInterval),
		gocron.NewTask(scheduler.runCleanUpTags),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)

	return scheduler, nil
}

func (s *Scheduler) Run() <-chan error {
	s.client.Start()
	return s.lifeErrs
}

func (s *Scheduler) Shutdown(ctx context.Context) error {
	shutdownDone := make(chan struct{})
	go func() {
		s.client.Shutdown()
		close(shutdownDone)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-shutdownDone:
	}
	return nil
}

func (s *Scheduler) runCleanUpTags() {
	defer func() {
		if v := recover(); v != nil {
			slog.Error("cron goroutine panicked", "recover", v, "stack", string(debug.Stack()))
			s.lifeErrs <- fmt.Errorf("panic during tag cleanup: %v", v)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx, span := tracing.StartSpan(ctx, "cron.Scheduler.runCleanUpTags")
	defer span.End()

	deleteCount, err := s.memoService.DeleteOrphanTags(ctx)
	if err != nil {
		slog.Error("failed to delete orphan tags", "error", err)
		return
	}
	if deleteCount > 0 {
		slog.Info("deleted orphan tags", "count", deleteCount)
	}
}
