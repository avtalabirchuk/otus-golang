package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	limit         int64
	offset        int64
	checkFileName string
}

func TestCopy(t *testing.T) {
	for _, tst := range []test{
		{
			limit:         0,
			offset:        0,
			checkFileName: "testdata/out_offset0_limit0.txt",
		},
	} {
		t.Run(fmt.Sprintf("Limit: %d, Offset: %d", tst.limit, tst.offset), func(t *testing.T) {
			targetFile, err := ioutil.TempFile("/tmp/", "output")
			require.NoError(t, err)
			defer os.Remove(targetFile.Name())

			err = Copy("testdata/input.txt", targetFile.Name(), tst.offset, tst.limit)
			require.NoError(t, err)

			checkFile, _ := os.Open(tst.checkFileName)
			defer checkFile.Close()

			checkFileContent, err := ioutil.ReadAll(checkFile)
			if err != nil {
				log.Fatal(err)
			}
			targetFileContent, err := ioutil.ReadAll(targetFile)
			if err != nil {
				log.Fatal(err)
			}
			require.Equal(t, checkFileContent, targetFileContent)
		})
	}
}
