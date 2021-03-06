package nanny

import (
	"context"
	"fmt"
	"runtime"
)

const (
	defaultStackSize = 4 << 10
)

// WithRecovery returns a middleware which recovers from panics.
func WithRecovery() RouteOptionFn {
	return func(r *route) {
		m := func(next Handler) Handler {
			return func(ctx context.Context, req Request) (resp interface{}, err error) {
				defer func() {
					if rec := recover(); rec != nil {
						errFromPanic, ok := rec.(error)
						if !ok {
							errFromPanic = fmt.Errorf("%v", rec)
						}
						stack := make([]byte, defaultStackSize)
						length := runtime.Stack(stack, false)
						msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", errFromPanic, stack[:length])

						r.logger.Println(msg)
						err = panicErr
					}
				}()

				resp, err = next(ctx, req)
				return
			}
		}

		r.middlewares = append(r.middlewares, m)
	}
}
