package worker

import (
	"context"
	"log"
	"time"

	"canceler/internal/canceler"

	"github.com/go-co-op/gocron"
)

type Worker struct {
	jobCanceler *canceler.JobCanceler

	scheduler *gocron.Scheduler
}

func New(cancellationTimeout time.Duration) *Worker {
	scheduler := gocron.NewScheduler(time.UTC)
	jobCanceler := canceler.New(scheduler.Stop, cancellationTimeout)
	jobCanceler.OnNormalCancel(func() {
		log.Println("normal cancel")
	})
	jobCanceler.OnTimeoutCancel(func() {
		log.Println("timeout cancel")
	})

	return &Worker{
		jobCanceler: jobCanceler,
		scheduler:   scheduler,
	}
}

func (w *Worker) ScheduleRunner(runner func(context.Context), interval time.Duration) {
	job, _ := w.scheduler.Every(interval).
		Do(runner, w.jobCanceler.CancellationContext())

	job.SingletonMode()
}

func (w *Worker) ScheduleServiceRepeatableRunner(
	runner func(context.Context) bool,
	interval time.Duration,
) {
	w.ScheduleRunner(func(ctx context.Context) {
		repeat := true

		for repeat {
			select {
			case <-ctx.Done():
				return
			default:
				repeat = runner(ctx)
			}
		}
	}, interval)
}

func (w *Worker) StartAsync() {
	w.scheduler.StartAsync()
}

func (w *Worker) Stop() {
	w.jobCanceler.Cancel()
}
