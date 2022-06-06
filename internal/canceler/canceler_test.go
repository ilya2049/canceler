package canceler_test

import (
	"canceler/internal/canceler"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJobCanceler_Cancel_CancelJobCalled(t *testing.T) {
	cancelJobCalled := false

	aCanceler := canceler.New(func() {
		cancelJobCalled = true
	}, time.Hour)

	aCanceler.Cancel()

	assert.True(t, cancelJobCalled)
}

func TestJobCanceler_Cancel_NormalCancel(t *testing.T) {
	isNormalCancel := false

	aCanceler := canceler.New(func() {}, time.Hour).
		OnNormalCancel(func() {
			isNormalCancel = true
		})

	aCanceler.Cancel()

	assert.True(t, isNormalCancel)
}

func TestJobCanceler_Cancel_TimeoutCancel(t *testing.T) {
	isTimeoutCancel := false

	aCanceler := canceler.New(func() {
		time.Sleep(time.Hour)
	}, time.Nanosecond).
		OnTimeoutCancel(func() {
			isTimeoutCancel = true
		})

	aCanceler.Cancel()

	assert.True(t, isTimeoutCancel)
}

func TestJobCanceler_CancellationContext_Done(_ *testing.T) {
	aCanceler := canceler.New(func() {}, time.Hour)

	cancellationContext := aCanceler.CancellationContext()
	go func() {
		aCanceler.Cancel()
	}()

	<-cancellationContext.Done()
}
