package main

import (
	"io"
)

type FakeCloseReader struct {
	io.Reader
}

func NewFakeCloseReader(r io.Reader) *FakeCloseReader {
	if r == io.Reader(nil) {
		return nil
	}
	return &FakeCloseReader{
		r,
	}
}

func (f *FakeCloseReader) Close() error {
	return nil
}

func (f *FakeCloseReader) Read(p []byte) (n int, err error) {
	if f == nil {
		return 0, io.EOF
	}
	return f.Reader.Read(p)
}
