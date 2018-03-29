package id3v2

import (
	"bufio"
	"io"
)

type bufWriter struct {
	err     error
	w       *bufio.Writer
	written int64
}

func newBufWriter(w io.Writer) *bufWriter {
	return &bufWriter{w: bufio.NewWriter(w)}
}

func (bw *bufWriter) EncodeAndWriteText(src string, to Encoding) {
	if bw.err != nil {
		return
	}
	bw.err = encodeWriteText(bw, src, to)
}

func (bw *bufWriter) Flush() error {
	if bw.err != nil {
		return bw.err
	}
	return bw.w.Flush()
}

func (bw *bufWriter) Reset(w io.Writer) {
	bw.err = nil
	bw.written = 0
	bw.w.Reset(w)
}

func (bw *bufWriter) WriteByte(c byte) {
	if bw.err != nil {
		return
	}
	bw.err = bw.w.WriteByte(c)
	if bw.err == nil {
		bw.written += 1
	}
}

func (bw *bufWriter) WriteBytesSize(size uint) {
	if bw.err != nil {
		return
	}
	bw.err = writeBytesSize(bw, size)
}

func (bw *bufWriter) WriteString(s string) {
	if bw.err != nil {
		return
	}

	var n int
	n, bw.err = bw.w.WriteString(s)
	bw.written += int64(n)
}

func (bw *bufWriter) Write(p []byte) (n int, err error) {
	if bw.err != nil {
		return 0, bw.err
	}
	n, err = bw.w.Write(p)
	bw.written += int64(n)
	bw.err = err
	return n, err
}

func (bw *bufWriter) Written() int64 {
	return bw.written
}

func useBufWriter(w io.Writer, f func(*bufWriter)) (int64, error) {
	var writtenBefore int64
	bw, ok := w.(*bufWriter)
	if ok {
		writtenBefore = bw.Written()
	} else {
		bw = getBufWriter(w)
		defer putBufWriter(bw)
	}

	f(bw)

	return bw.Written() - writtenBefore, bw.Flush()
}