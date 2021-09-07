package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		env, err := ReadDir("testdata/tmp")
		require.NoError(t, err)
		expectedEnv := Environment{
			"BAR": "bar",
		}
		require.Equal(t, expectedEnv, env)
	})
}
