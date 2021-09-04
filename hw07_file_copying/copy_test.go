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
		{
			limit:         10,
			offset:        0,
			checkFileName: "testdata/out_offset0_limit10.txt",
		},
		{
			limit:         1000,
			offset:        0,
			checkFileName: "testdata/out_offset0_limit1000.txt",
		},
		{
			limit:         10000,
			offset:        0,
			checkFileName: "testdata/out_offset0_limit10000.txt",
		},
		{
			limit:         1000,
			offset:        100,
			checkFileName: "testdata/out_offset100_limit1000.txt",
		},
		{
			limit:         10000,
			offset:        6000,
			checkFileName: "testdata/out_offset6000_limit1000.txt",
		},
	} {
		t.Run(fmt.Sprintf("Limit: %d, Offset: %d", tst.limit, tst.offset), func(t *testing.T) {
			targetFile, err := os.CreateTemp("/tmp/", "output")
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

		t.Run("if offset bigger than file size", func(t *testing.T) {
			err := Copy("testdata/input.txt", "/tmp/output.txt", 999999, 0)
			require.Equal(t, err, ErrOffsetExceedsFileSize)
		})

		t.Run("Copy from non file", func(t *testing.T) {
			err := Copy("testdata", "/tmp/output.txt", 123, 0)
			require.Equal(t, err, ErrUnsupportedFile)
		})
		t.Run("Copy to non exist directory", func(t *testing.T) {
			err := Copy("testdata/input.txt", "/tmp/non-exitst-dir/output.txt", 123, 0)
			require.Error(t, err)
		})

		t.Run("Copy from non exist file", func(t *testing.T) {
			err := Copy("testdata/non_exist_file", "/tmp/output.txt", 123, 0)
			require.Error(t, err)
		})
	}
}
