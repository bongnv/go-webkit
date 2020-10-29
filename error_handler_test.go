package webkit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WithErrorHandler(t *testing.T) {
	opt := WithErrorHandler(defaultErrorHandler(nil))
	r := &route{}
	opt(r)
	require.NotNil(t, r.errorHandler)
}
