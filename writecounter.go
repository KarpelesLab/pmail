package pmail

import "io"

type writeCounter struct {
	W io.Writer
	C int64
}

func (w *writeCounter) Write(b []byte) (int, error) {
	n, err := w.W.Write(b)
	w.C += int64(n)
	return n, err
}
