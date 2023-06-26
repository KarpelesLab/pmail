package pmail

import "io"

type lineBreaker struct {
	line []byte
	eol  []byte
	used int
	out  io.Writer
}

func newStdLinebreaker(w io.Writer) io.WriteCloser {
	return &lineBreaker{line: make([]byte, 78), eol: []byte{'\r', '\n'}}
}

func (l *lineBreaker) Write(b []byte) (n int, err error) {
	if l.used+len(b) < len(l.line) {
		copy(l.line[l.used:], b)
		l.used += len(b)
		return len(b), nil
	}

	n, err = l.out.Write(l.line[0:l.used])
	if err != nil {
		return
	}
	excess := len(l.line) - l.used
	l.used = 0

	n, err = l.out.Write(b[0:excess])
	if err != nil {
		return
	}

	n, err = l.out.Write(l.eol)
	if err != nil {
		return
	}

	return l.Write(b[excess:])
}

func (l *lineBreaker) Close() (err error) {
	if l.used > 0 {
		_, err = l.out.Write(l.line[0:l.used])
		if err != nil {
			return
		}
		_, err = l.out.Write(l.eol)
	}

	return
}
