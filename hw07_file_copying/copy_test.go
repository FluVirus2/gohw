package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SameFiles(f1, f2 *os.File) (bool, error) {
	buf1 := make([]byte, 4096)
	buf2 := make([]byte, 4096)

	for {
		n1, err1 := f1.Read(buf1)
		if err1 != nil && !errors.Is(err1, io.EOF) {
			return false, fmt.Errorf("ошибка чтения из первого файла: %w", err1)
		}

		n2, err2 := f2.Read(buf2)
		if err2 != nil && !errors.Is(err2, io.EOF) {
			return false, fmt.Errorf("ошибка чтения из второго файла: %w", err2)
		}

		if n1 != n2 {
			return false, nil
		}

		for i := 0; i < n1; i++ {
			if buf1[i] != buf2[i] {
				return false, nil
			}
		}

		if errors.Is(err1, io.EOF) && errors.Is(err2, io.EOF) {
			return true, nil
		}

		if (errors.Is(err1, io.EOF) && !errors.Is(err2, io.EOF)) || (!errors.Is(err1, io.EOF) && errors.Is(err2, io.EOF)) {
			return false, nil
		}
	}
}

func TestCopyOk(t *testing.T) {
	type TestCase struct {
		Name string

		InputPath  string
		OutputPath string
		Limit      int64
		Offset     int64

		ExpectedPath string
	}

	testCases := []TestCase{
		{
			Name:       "offet 0 limit 0",
			InputPath:  "testdata/input.txt",
			OutputPath: "testdata/output1.txt",
			Limit:      0,
			Offset:     0,

			ExpectedPath: "testdata/out_offset0_limit0.txt",
		},
		{
			Name:       "offet 0 limit 10",
			InputPath:  "testdata/input.txt",
			OutputPath: "testdata/output2.txt",
			Limit:      10,
			Offset:     0,

			ExpectedPath: "testdata/out_offset0_limit10.txt",
		},
		{
			Name:       "offet 0 limit 1000",
			InputPath:  "testdata/input.txt",
			OutputPath: "testdata/output3.txt",
			Limit:      1000,
			Offset:     0,

			ExpectedPath: "testdata/out_offset0_limit1000.txt",
		},
		{
			Name:       "offet 0 limit 10000",
			InputPath:  "testdata/input.txt",
			OutputPath: "testdata/output4.txt",
			Limit:      10000,
			Offset:     0,

			ExpectedPath: "testdata/out_offset0_limit10000.txt",
		},
		{
			Name:       "offet 100 limit 1000",
			InputPath:  "testdata/input.txt",
			OutputPath: "testdata/output5.txt",
			Limit:      1000,
			Offset:     100,

			ExpectedPath: "testdata/out_offset100_limit1000.txt",
		},
		{
			Name:       "offet 6000 limit 1000",
			InputPath:  "testdata/input.txt",
			OutputPath: "testdata/output6.txt",
			Limit:      1000,
			Offset:     6000,

			ExpectedPath: "testdata/out_offset100_limit1000.txt",
		},
		{
			Name:       "empty file",
			InputPath:  "testdata/.gitkeep",
			OutputPath: "testdata/output7.txt",
			Limit:      0,
			Offset:     0,

			ExpectedPath: "testdata/.gitkeep",
		},
	}

	t.Parallel()
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := Copy(testCase.InputPath, testCase.OutputPath, testCase.Offset, testCase.Limit)
			defer os.Remove(testCase.OutputPath)

			assert.NoError(t, err)
			assert.FileExists(t, testCase.OutputPath)

			product, err := os.Open(testCase.OutputPath)
			assert.NoError(t, err)

			expected, err := os.Open(testCase.ExpectedPath)
			assert.NoError(t, err)

			res, err := SameFiles(product, expected)
			assert.NoError(t, err)
			assert.True(t, res)

			err = product.Close()
			assert.NoError(t, err)

			err = expected.Close()
			assert.NoError(t, err)
		})
	}
}

func TestCopyError(t *testing.T) {
	type TestCase struct {
		Name          string
		InputPath     string
		OutputPath    string
		Limit         int64
		Offset        int64
		ExpectedError error
	}

	testCases := []TestCase{
		{
			Name:          "offset more than file size",
			InputPath:     "testdata/input.txt",
			OutputPath:    "testdata/output8.txt",
			Limit:         0,
			Offset:        1 << 60,
			ExpectedError: ErrOffsetExceedsFileSize,
		},
		{
			Name:          "unsupported file type",
			InputPath:     "testdata",
			OutputPath:    "testdata/output9.txt",
			Limit:         0,
			Offset:        0,
			ExpectedError: ErrUnsupportedFile,
		},
	}

	t.Parallel()
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := Copy(testCase.InputPath, testCase.OutputPath, testCase.Offset, testCase.Limit)
			defer os.Remove(testCase.OutputPath)

			assert.ErrorIs(t, err, testCase.ExpectedError)
		})
	}
}
