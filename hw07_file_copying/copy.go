package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrInvalidParam          = errors.New("'-offset' and '-limit' must be >= 0; '-from' and '-to' must not be empty")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 || limit < 0 || fromPath == "" || toPath == "" {
		return ErrInvalidParam
	}

	srcF, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer closeWithCheck(srcF)

	srcFi, err := srcF.Stat()
	if err != nil {
		return err
	}

	if srcFi.Size() == 0 {
		return ErrUnsupportedFile
	}

	if offset > srcFi.Size() {
		return ErrOffsetExceedsFileSize
	}

	if _, err = srcF.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	dstF, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer closeWithCheck(dstF)

	if limit == 0 || limit > srcFi.Size()-offset {
		limit = srcFi.Size() - offset
	}

	bar := pb.Full.Start64(limit)
	defer bar.Finish()
	_, err = io.CopyN(dstF, bar.NewProxyReader(srcF), limit)
	return err
}

func closeWithCheck(f *os.File) {
	if clErr := f.Close(); clErr != nil {
		fmt.Printf("problem during closing %s; details: %s\n", f.Name(), clErr.Error())
	}
}
