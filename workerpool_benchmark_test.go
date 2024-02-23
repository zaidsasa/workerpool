package workerpool_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/zaidsasa/workerpool"
)

const (
	poolSize      = 5e4
	runCount      = 1e6
	sleepDuration = 10 * time.Millisecond
)

func stubFunc() {
	time.Sleep(sleepDuration)
}

func BenchmarkGoroutines(b *testing.B) {
	var waitGroup sync.WaitGroup

	for i := 0; i < b.N; i++ {
		waitGroup.Add(runCount)

		for j := 0; j < runCount; j++ {
			go func() {
				stubFunc()
				waitGroup.Done()
			}()
		}
		waitGroup.Wait()
	}
}

func BenchmarkWorkerPool(b *testing.B) {
	wPool, _ := workerpool.New(poolSize)

	for i := 0; i < b.N; i++ {
		for j := 0; j < runCount; j++ {
			wPool.Submit(stubFunc)
		}
	}

	_ = wPool.Wait(context.Background())
}
