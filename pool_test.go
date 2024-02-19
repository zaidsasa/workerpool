package workerpool_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zaidsasa/workerpool"
)

func TestNewPool(t *testing.T) {
	t.Parallel()

	type args struct {
		size int
		opts []workerpool.Opt
	}

	testCases := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "pool with size of 0",
			args: args{
				size: 0,
				opts: nil,
			},
			wantErr: workerpool.ErrInvalidPoolSize,
		},
		{
			name: "pool with size of 10",
			args: args{
				size: 3,
				opts: nil,
			},
			wantErr: nil,
		},
		{
			name: "pool with size of 100",
			args: args{
				size: 100,
				opts: nil,
			},
			wantErr: nil,
		},
		{
			name: "pool with size of 100 and WithKeepAliveOption",
			args: args{
				size: 100,
				opts: []workerpool.Opt{workerpool.WithKeepAliveOption()},
			},
			wantErr: nil,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			gotPool, err := workerpool.NewPool(test.args.size, test.args.opts...)
			if test.wantErr != nil {
				assert.ErrorIs(t, err, test.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotPool)
			}
		})
	}
}

func TestPool_Wait(t *testing.T) {
	t.Parallel()

	type fields struct {
		PoolOpts []workerpool.Opt
	}

	type args struct {
		ctxFunc func() (context.Context, context.CancelFunc)
	}

	testCases := []struct {
		name              string
		fields            fields
		args              args
		wantErr           error
		wantWaitErr       error
		forceCancellation bool
	}{
		{
			name: "wait pool with default behavior",
			fields: fields{
				PoolOpts: nil,
			},
			args: args{
				ctxFunc: func() (context.Context, context.CancelFunc) {
					return context.Background(), nil
				},
			},
		},
		{
			name: "wait pool with timeout context and keep alive option",
			fields: fields{
				PoolOpts: []workerpool.Opt{workerpool.WithKeepAliveOption()},
			},
			args: args{
				ctxFunc: func() (context.Context, context.CancelFunc) {
					return context.WithTimeout(context.Background(), 1*time.Second)
				},
			},
			wantWaitErr: context.DeadlineExceeded,
		},
		{
			name: "wait pool with cancel context",
			fields: fields{
				PoolOpts: nil,
			},
			args: args{
				ctxFunc: func() (context.Context, context.CancelFunc) {
					return context.WithCancel(context.Background())
				},
			},
			forceCancellation: true,
			wantWaitErr:       context.Canceled,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			pool, err := workerpool.NewPool(10, test.fields.PoolOpts...)
			if test.wantErr != nil {
				assert.ErrorIs(t, err, test.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pool)
			}

			ctx, cancel := test.args.ctxFunc()
			if cancel != nil {
				if test.forceCancellation {
					cancel()
				} else {
					defer cancel()
				}
			}

			requests := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
			responseChan := make(chan string, len(requests))
			for _, r := range requests {
				r := r
				pool.Run(func() {
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

			if !test.forceCancellation {
				assert.ElementsMatch(t, expectedResp, requests)
			}
		})
	}
}
