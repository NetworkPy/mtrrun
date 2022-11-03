package agent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCounter(t *testing.T) {
	name, help := "test", "help"

	c := NewCounter(name, help)
	require.NotNil(t, c)

	d := c.Desc()
	require.Equal(t, name, d.Name)
	require.Equal(t, help, d.Help)
}
