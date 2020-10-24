package webkit

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handler defines a function to serve HTTP requests.
type Handler func(ctx context.Context, req Request) error

// Middleware defines a middleware to provide additional logic.
type Middleware func(Handler) Handler

func buildHandlerFunc(h Handler) httprouter.Handle {
	return func(w http.ResponseWriter, httpReq *http.Request, params httprouter.Params) {
		ctx := httpReq.Context()
		req := &requestImpl{
			httpWriter: w,
			httpReq:    httpReq,
		}

		err := h(ctx, req)
		if err != nil {
			req.responseError(err)
		}
	}
}
