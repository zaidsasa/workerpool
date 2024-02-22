package workerpool_test

import (
	"context"
	"testing"

	"github.com/zaidsasa/workerpool"
)

const poolSize = 10

func stubFunc() {}

func BenchmarkWorkerPool(b *testing.B) {
	wPool, _ := workerpool.New(poolSize)

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			wPool.Submit(stubFunc)
		}
	}

	_ = wPool.Wait(context.Background())
}