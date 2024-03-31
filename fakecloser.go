package main

import "io"

type FakeCloseReader struct {
	reader io.Reader
}

func NewFakeCloseReader(r io.Reader) FakeCloseReader {
	return FakeCloseReader{
		reader: r,
	}
}

func (f FakeCloseReader) Close() error {
	return nil
}

func (f FakeCloseReader) Read(p []byte) (n int, err error) {
	return f.reader.Read(p)
}
