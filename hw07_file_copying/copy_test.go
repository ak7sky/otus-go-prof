package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/udhos/equalfile"
)

func TestCopy(t *testing.T) {
	testCases := []struct {
		src, dst, exp string
		ofs, lim      int64
	}{
		{src: "testdata/input.txt", dst: "out.txt", ofs: 0, lim: 0, exp: "testdata/out_offset0_limit0.txt"},
		{src: "testdata/input.txt", dst: "out.txt", ofs: 0, lim: 10, exp: "testdata/out_offset0_limit10.txt"},
		{src: "testdata/input.txt", dst: "out.txt", ofs: 0, lim: 1000, exp: "testdata/out_offset0_limit1000.txt"},
		{src: "testdata/input.txt", dst: "out.txt", ofs: 0, lim: 10000, exp: "testdata/out_offset0_limit10000.txt"},
		{src: "testdata/input.txt", dst: "out.txt", ofs: 100, lim: 1000, exp: "testdata/out_offset100_limit1000.txt"},
		{src: "testdata/input.txt", dst: "out.txt", ofs: 6000, lim: 1000, exp: "testdata/out_offset6000_limit1000.txt"},
	}

	for _, tc := range testCases {
		tc := tc
		tcName := fmt.Sprintf("%s-->%s(offset:%d;limit:%d)", tc.src, tc.dst, tc.ofs, tc.lim)
		t.Run(tcName, func(t *testing.T) {
			err := Copy(tc.src, tc.dst, tc.ofs, tc.lim)
			require.NoError(t, err)

			cmp := equalfile.New(nil, equalfile.Options{})
			equal, err := cmp.CompareFile(tc.dst, tc.exp)
			if err != nil {
				require.Fail(t, "err during comparison actual with expected result")
			}
			require.True(t, equal)
			if err := os.Remove(tc.dst); err != nil {
				require.Fail(t, "err during removing dst file after test")
			}
		})
	}
}

func TestCopyFailure(t *testing.T) {
	testCases := []struct {
		src, dst, expErrMsg string
		ofs, lim            int64
	}{
		{src: "testdata/input.txt", dst: "out.txt", ofs: -1, lim: 0, expErrMsg: ErrInvalidParam.Error()},
		{src: "testdata/input.txt", dst: "out.txt", ofs: 0, lim: -1, expErrMsg: ErrInvalidParam.Error()},
		{src: "testdata/input.txt", dst: "out.txt", ofs: -1, lim: -1, expErrMsg: ErrInvalidParam.Error()},
		{src: "testdata/input.txt", dst: "out.txt", ofs: 1000_000, lim: 0, expErrMsg: ErrOffsetExceedsFileSize.Error()},
		{src: "testdata/non-existent.txt", dst: "out.txt", ofs: 0, lim: 0, expErrMsg: "no such file or directory"},
		{src: "testdata/empty_input", dst: "out.txt", ofs: 0, lim: 0, expErrMsg: ErrUnsupportedFile.Error()},
	}

	for _, tc := range testCases {
		tc := tc
		tcName := fmt.Sprintf("%s-->%s(offset:%d;limit:%d)", tc.src, tc.dst, tc.ofs, tc.lim)
		t.Run(tcName, func(t *testing.T) {
			err := Copy(tc.src, tc.dst, tc.ofs, tc.lim)
			require.Error(t, err)
			require.ErrorContains(t, err, tc.expErrMsg)
		})
	}
}
