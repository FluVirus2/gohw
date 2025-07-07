package hw02unpackstring

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	type TestCase struct {
		Name     string
		Input    string
		Expected string
	}

	tests := []TestCase{
		{
			Name:     "a4bc2d5e",
			Input:    "a4bc2d5e",
			Expected: "aaaabccddddde",
		},
		{
			Name:     "abccd",
			Input:    "abccd",
			Expected: "abccd",
		},
		{
			Name:     "empty string",
			Input:    "",
			Expected: "",
		},
		{
			Name:     "aaa0b",
			Input:    "aaa0b",
			Expected: "aab",
		},
		{
			Name:     "🙃0",
			Input:    "🙃0",
			Expected: "",
		},
		{
			Name:     "aaф0b",
			Input:    "aaф0b",
			Expected: "aab",
		},
		{
			Name:     "a2b2",
			Input:    "a2b2",
			Expected: "aabb",
		},
		{
			Name:     "a2b",
			Input:    "a2b",
			Expected: "aab",
		},
		{
			Name:     "ab2",
			Input:    "ab2",
			Expected: "abb",
		},
		{
			Name:     "null characters 1",
			Input:    "\u00004",
			Expected: "\u0000\u0000\u0000\u0000",
		},
		{
			Name:     "null characters 2",
			Input:    "a\u00004",
			Expected: "a\u0000\u0000\u0000\u0000",
		},
		{
			Name:     "null characters 3",
			Input:    "\u00004a",
			Expected: "\u0000\u0000\u0000\u0000a",
		},
		{
			Name:     "null characters 4",
			Input:    "a\u00004a",
			Expected: "a\u0000\u0000\u0000\u0000a",
		},
		{
			Name:     "nonprintables",
			Input:    "\u00022\u00163\u00194",
			Expected: "\u0002\u0002\u0016\u0016\u0016\u0019\u0019\u0019\u0019",
		},
		{
			Name:     "hieroglyphs",
			Input:    "\u99782\u73483\u8345",
			Expected: "饸饸獈獈獈荅",
		},
		{
			Name:     "large string",
			Input:    strings.Repeat("a2", 1_000_000) + strings.Repeat("b3", 1_000_000),
			Expected: strings.Repeat("aa", 1_000_000) + strings.Repeat("bbb", 1_000_000),
		},
		// uncomment if task with asterisk completed
		// {Input: `qwe\4\5`, Expected: `qwe45`},
		// {Input: `qwe\45`, Expected: `qwe44444`},
		// {Input: `qwe\\5`, Expected: `qwe\\\\\`},
		// {Input: `qwe\\\3`, Expected: `qwe\3`},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := Unpack(tc.Input)
			require.NoError(t, err)
			require.Equal(t, tc.Expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	type TestCase struct {
		Name          string
		Input         string
		ExpectedError error
	}

	testCases := []TestCase{
		{
			Name:          "3abc",
			Input:         "3abc",
			ExpectedError: ErrInvalidString,
		},
		{
			Name:          "45",
			Input:         "45",
			ExpectedError: ErrInvalidString,
		},
		{
			Name:          "aaa10b",
			Input:         "aaa10b",
			ExpectedError: ErrInvalidString,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := Unpack(tc.Input)
			require.Truef(t, errors.Is(err, tc.ExpectedError), "actual error: %q", err)
		})
	}
}
