package gwf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Plugin_Apply(t *testing.T) {
	app := New()
	DefaultApp.Apply(app)
	require.Len(t, app.routeOptions, 5)
	require.NotNil(t, app.logger)
}