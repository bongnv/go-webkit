package nanny

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Plugin_Apply(t *testing.T) {
	app := New()
	DefaultApp.Apply(app)
	require.Len(t, app.routeOptions, 8)
	require.NotNil(t, app.logger)
}
