package nanny

import (
	"context"
	"time"
)

// WithTimeout specifies the time limit  for a route.
func WithTimeout(timeout time.Duration) RouteOptionFn {
	return func(r *route) {
		r.timeout = timeout
	}
}

func injectTimeoutMiddleware() RouteOptionFn {
	return func(r *route) {
		var m Middleware = func(next Handler) Handler {
			if r.timeout == 0 {
				return next
			}

			return func(ctx context.Context, req Request) (interface{}, error) {
				resultCh := make(chan handlerResult, 1)
				go func() {
					defer func() {
						if r := recover(); r != nil {
							// TODO: Add logs
							resultCh <- handlerResult{err: panicErr}
						}
						close(resultCh)
					}()

					resp, err := next(ctx, req)
					resultCh <- handlerResult{resp: resp, err: err}
				}()

				timeCh := time.NewTimer(r.timeout)
				defer timeCh.Stop()
				select {
				case r := <-resultCh:
					return r.resp, r.err
				case <-timeCh.C:
					return nil, timeoutErr
				}
			}
		}

		r.middlewares = append(r.middlewares, m)
	}
}

type handlerResult struct {
	resp interface{}
	err  error
}
