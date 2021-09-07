package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
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

	t.Run("Directory does't exist", func(t *testing.T) {
		_, err := ReadDir("non_exist_dir")
		require.EqualError(t, err, "directory does't exist")
	})

}
