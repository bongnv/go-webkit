package gwf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Plugin_Apply(t *testing.T) {
	app := New(DefaultApp)
	require.Len(t, app.routeOptions, 5)
	require.NotNil(t, app.router)
	require.NotNil(t, app.logger)
}
