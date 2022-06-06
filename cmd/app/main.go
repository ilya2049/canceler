package main

import (
	"canceler/internal/worker"
	"context"
	"fmt"

	"time"
)

func main() {
	aWorker := worker.New(1 * time.Second)

	var i int
	aWorker.ScheduleServiceRepeatableRunner(func(_ context.Context) bool {
		fmt.Println(i)
		i++

		return i < 5
	}, time.Second)
	aWorker.StartAsync()

	time.Sleep(4 * time.Second)
	aWorker.Stop()
}
