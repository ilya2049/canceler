package canceler

import (
	"context"
	"time"
)

type JobCanceler struct {
	cancelJob func()

	cancellationTimeout time.Duration
	cancellationCtx     context.Context
	startCancellation   context.CancelFunc

	onNormalCancelFunc  func()
	onTimeoutCancelFunc func()
}

func New(cancelJob func(), cancellationTimeout time.Duration) *JobCanceler {
	ctx, cancel := context.WithCancel(context.Background())

	return &JobCanceler{
		cancelJob:           cancelJob,
		cancellationTimeout: cancellationTimeout,
		cancellationCtx:     ctx,
		startCancellation:   cancel,

		onNormalCancelFunc:  func() {},
		onTimeoutCancelFunc: func() {},
	}
}

func (c *JobCanceler) OnNormalCancel(f func()) *JobCanceler {
	c.onNormalCancelFunc = f

	return c
}

func (c *JobCanceler) OnTimeoutCancel(f func()) *JobCanceler {
	c.onTimeoutCancelFunc = f

	return c
}

func (c *JobCanceler) CancellationContext() context.Context {
	return c.cancellationCtx
}

func (c *JobCanceler) Cancel() {
	ctx, cancel := context.WithTimeout(context.Background(), c.cancellationTimeout)
	defer cancel()

	c.startCancellation()

	cancelJobChan := make(chan struct{})
	go func() {
		c.cancelJob()
		close(cancelJobChan)
	}()

	select {
	case <-ctx.Done():
		c.onTimeoutCancelFunc()

	case <-cancelJobChan:
		c.onNormalCancelFunc()
	}
}
