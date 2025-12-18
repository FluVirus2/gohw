package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func IndicateProgress(progress float64) {
	fmt.Printf("Copied %.2f%%", progress*100)
	fmt.Println()
}

func IndicateStart() {
	fmt.Println("Copying is staring")
}

func IndicateEnd() {
	fmt.Println("Copying is finished")
}

type ProgressWriter struct {
	Writer     io.Writer
	Written    int64
	Total      int64
	OnProgress func(float64)
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n, err := pw.Writer.Write(p)
	pw.Written += int64(n)
	if pw.OnProgress != nil {
		progress := float64(pw.Written) / float64(pw.Total)
		pw.OnProgress(progress)
	}

	return n, err
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	inputFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("cannot open source file: %w", err)
	}
	defer inputFile.Close()

	stats, err := inputFile.Stat()
	if err != nil {
		return fmt.Errorf("cannot get file stats: %w", err)
	}

	if !stats.Mode().IsRegular() || stats.IsDir() {
		return ErrUnsupportedFile
	}

	size := stats.Size()
	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		limit = size - offset
	}

	_, err = inputFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("cannot seek to offset: %w", err)
	}

	outputFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("cannot create destination file: %w", err)
	}
	defer outputFile.Close()

	progressWriter := &ProgressWriter{
		Writer:     outputFile,
		Written:    0,
		Total:      limit,
		OnProgress: IndicateProgress,
	}

	IndicateStart()
	_, err = io.CopyN(progressWriter, inputFile, limit)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error during copy: %w", err)
	}

	IndicateEnd()
	return nil
}
