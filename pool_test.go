package workerpool

import (
	"context"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	p, err := NewPool(6)
	if err != nil {
		t.Errorf("NewPool() error = %v, wantErr %v", err, true)
		return
	}

	requests := []string{"alpha", "beta", "gamma", "delta", "epsilon"}

	rspChan := make(chan string, len(requests))
	for _, r := range requests {
		r := r
		p.Run(func() {
			rspChan <- r
		})
	}

	if err := p.Wait(context.Background()); err != nil {
		t.Fatal("Error should be empty")
	}

	close(rspChan)
	rspSet := map[string]struct{}{}
	for rsp := range rspChan {
		rspSet[rsp] = struct{}{}
	}
	if len(rspSet) < len(requests) {
		t.Fatal("Did not handle all requests")
	}
	for _, req := range requests {
		if _, ok := rspSet[req]; !ok {
			t.Fatal("Missing expected values:", req)
		}
	}
}

func TestNewPool_WithKeepAliveForTwoSeconds(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	p, err := NewPool(2, WithKeepAliveOption())
	if err != nil {
		t.Errorf("NewPool() error = %v, wantErr %v", err, true)
		return
	}

	requests := []string{"alpha", "beta", "gamma", "delta", "epsilon"}

	rspChan := make(chan string, len(requests))
	for _, r := range requests {
		r := r
		p.Run(func() {
			rspChan <- r
		})
	}

	if err := p.Wait(ctx); err != context.DeadlineExceeded {
		t.Fatalf("Error should be of type %q", context.DeadlineExceeded)
	}

	close(rspChan)
	rspSet := map[string]struct{}{}
	for rsp := range rspChan {
		rspSet[rsp] = struct{}{}
	}
	if len(rspSet) < len(requests) {
		t.Fatal("Did not handle all requests")
	}
	for _, req := range requests {
		if _, ok := rspSet[req]; !ok {
			t.Fatal("Missing expected values:", req)
		}
	}
}
