package countingWriter

import "io"

type countingWriter struct {
	w io.Writer
	c int64
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	cw.c += int64(len(p))
	return cw.w.Write(p)
}

// CountingWriter returns a wrapped up io.Writer and
// a pointer to a bytes counter
func CountingWriter(w io.Writer) (io.Writer, *int64) {
	cw := &countingWriter{w, 0}

	return cw, &cw.c
}
