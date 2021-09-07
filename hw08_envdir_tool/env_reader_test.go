package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		env, err := ReadDir("env/tmp")
		require.NoError(t, err)
		expectedEnv := Environment{
			"BAR":   "bar",
			"EMPTY": "",
			"FOO":   "   foo\nwith new line",
			"HELLO": `"hello"`,
			"UNSET": "",
		}
		require.Equal(t, expectedEnv, env)
	})

}
