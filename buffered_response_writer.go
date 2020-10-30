package gwf

import (
	"bytes"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func brwTransformer(next httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		bw := newBufRespWriter(rw)
		next(bw, req, params)
		_ = bw.Close()
	}
}

func newBufRespWriter(rw http.ResponseWriter) *bufRespWriter {
	bw := &bufRespWriter{
		Buffer:         &bytes.Buffer{},
		ResponseWriter: rw,
		statusCode:     http.StatusOK,
	}

	return bw
}

type bufRespWriter struct {
	*bytes.Buffer
	http.ResponseWriter
	statusCode int
}

func (bw *bufRespWriter) WriteHeader(code int) {
	bw.statusCode = code
}

func (bw bufRespWriter) Write(b []byte) (int, error) {
	if len(b) > 0 {
		return bw.Buffer.Write(b)
	}

	return 0, nil
}

func (bw bufRespWriter) Close() error {
	bw.ResponseWriter.WriteHeader(bw.statusCode)
	if bw.statusCode == http.StatusNoContent {
		return nil
	}

	_, err := bw.Buffer.WriteTo(bw.ResponseWriter)
	return err
}
