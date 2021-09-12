package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	// Place your code here
	t.Run("success", func(t *testing.T) {
		cmd := []string{"testdata/write_file.sh", "arg1=1", "arg2=2"}
		environment := Environment{
			"BAR":   "bar",
			"FOO":   "   foo\nwith new line",
			"HELLO": `"hello"`,
			"UNSET": "",
		}
		returnCode := RunCmd(cmd, environment)
		require.Equal(t, 0, returnCode)
	})
	t.Run("using env with empty key", func(t *testing.T) {
		environment := Environment{
			"": "value",
		}
		returnCode := RunCmd([]string{"command", "arg"}, environment)
		require.Equal(t, -1, returnCode)
	})
}
