package nanny

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_WithTimeout(t *testing.T) {
	r := &route{}
	WithTimeout(time.Second)(r)
	require.Equal(t, time.Second, r.timeout)
}

func Test_injectTimeoutMiddleware(t *testing.T) {
	opt := injectTimeoutMiddleware()

	t.Run("without-timeout", func(t *testing.T) {
		r := &route{}
		opt(r)
		require.Len(t, r.middlewares, 1)
		var h Handler = func(ctx context.Context, req Request) (interface{}, error) {
			panic("it will panic")
		}

		require.Panics(t, func() {
			r.middlewares[0](h)(context.Background(), nil)
		}, "The handler should panic as there no panic recovery")
	})

	t.Run("with-timeout", func(t *testing.T) {
		r := &route{timeout: time.Second}
		opt(r)
		require.Len(t, r.middlewares, 1)
		var h Handler = func(ctx context.Context, req Request) (interface{}, error) {
			panic("it will panic")
		}

		require.NotPanics(t, func() {
			_, err := r.middlewares[0](h)(context.Background(), nil)
			require.EqualError(t, err, "Service Unavailable")
		}, "The handler should not panic as there panic recovery in the middleware")
	})

	t.Run("handler-timeout", func(t *testing.T) {
		r := &route{timeout: 50 * time.Millisecond}
		opt(r)
		require.Len(t, r.middlewares, 1)
		var h Handler = func(ctx context.Context, req Request) (interface{}, error) {
			<-time.After(100 * time.Millisecond)
			return nil, nil
		}

		testDone := make(chan struct{})
		go func() {
			_, err := r.middlewares[0](h)(context.Background(), nil)
			require.EqualError(t, err, "Request Timeout")
			close(testDone)
		}()

		select {
		case <-time.After(75 * time.Millisecond):
			require.Fail(t, "Test timeout")
		case <-testDone:
		}
	})

	t.Run("handler-timeout", func(t *testing.T) {
		r := &route{timeout: 50 * time.Millisecond}
		opt(r)
		require.Len(t, r.middlewares, 1)
		var h Handler = func(ctx context.Context, req Request) (interface{}, error) {
			return nil, nil
		}

		testDone := make(chan struct{})
		go func() {
			_, err := r.middlewares[0](h)(context.Background(), nil)
			require.NoError(t, err)
			close(testDone)
		}()

		select {
		case <-time.After(75 * time.Millisecond):
			require.Fail(t, "Test timeout")
		case <-testDone:
		}
	})
}
