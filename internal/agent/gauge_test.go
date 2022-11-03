package agent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGauge(t *testing.T) {
	name, help := "test", "help"

	g := NewGauge(name, help)
	require.NotNil(t, g)

	d := g.Desc()
	require.Equal(t, name, d.Name)
	require.Equal(t, help, d.Help)
}
