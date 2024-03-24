package workerpool_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zaidsasa/workerpool"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		size    uint32
		options []workerpool.Option
	}

	testCases := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "workerpool creation failed with a size of 0.",
			args: args{
				size: 0,
			},
			wantErr: workerpool.ErrInvalidSize,
		},
		{
			name: "workerpool was created successfully with a size of 10",
			args: args{
				size: 10,
			},
			wantErr: nil,
		},
		{
			name: "workerpool was created successfully with a size of 100",
			args: args{
				size: 100,
			},
			wantErr: nil,
		},
		{
			name: "The worker pool was created successfully with a size of 100, and the 'WithKeepAliveOption' is enabled",
			args: args{
				size:    100,
				options: []workerpool.Option{workerpool.WithKeepAliveOption(true)},
			},
			wantErr: nil,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gotPool, err := workerpool.New(test.args.size, test.args.options...)

			if test.wantErr != nil {
				assert.ErrorIs(t, err, test.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotPool)
			}
		})
	}
}

func TestWorkerPool_Wait(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		WorkerPoolOptions  []workerpool.Option
		ctxFunc            func() (context.Context, context.CancelFunc)
		wantErr            error
		wantWaitErr        error
		submitCancellation bool
	}{
		{
			name: "workerpool terminated gracefully when the context of type 'background' was encountered",
			ctxFunc: func() (context.Context, context.CancelFunc) {
				return context.Background(), nil
			},
		},
		{
			name: "workerpool terminated gracefully upon receiving a context cancellation",
			ctxFunc: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			submitCancellation: true,
			wantWaitErr:        context.Canceled,
		},
		{
			name: "workerpool with the 'WithKeepAliveOption' option enabled terminated gracefully when the context timed out",

			WorkerPoolOptions: []workerpool.Option{workerpool.WithKeepAliveOption(true)},
			ctxFunc: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 2*time.Second)
			},
			wantWaitErr: context.DeadlineExceeded,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			pool, err := workerpool.New(10, test.WorkerPoolOptions...)

			if test.wantErr != nil {
				assert.ErrorIs(t, err, test.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pool)
			}

			ctx, cancel := test.ctxFunc()
			if cancel != nil {
				if test.submitCancellation {
					cancel()
				} else {
					defer cancel()
				}
			}

			requests := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
			responseChan := make(chan string, len(requests))

			for _, r := range requests {
				r := r

				pool.Submit(func() {
					time.Sleep(1 * time.Second)
					responseChan <- r
				})
			}

			err = pool.Wait(ctx)
			if test.wantWaitErr != nil {
				assert.ErrorIs(t, err, test.wantWaitErr)
			} else {
				assert.NoError(t, err)
			}

			close(responseChan)

			expectedResp := []string{}

			for rsp := range responseChan {
				expectedResp = append(expectedResp, rsp)
			}

			assert.ElementsMatch(t, expectedResp, requests)
		})
	}
}
