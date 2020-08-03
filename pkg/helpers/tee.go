package helpers

import (
	"io"

	"github.com/pkg/errors"
)

func TeeReadCloser(r io.ReadCloser, w io.Writer) io.ReadCloser {
	return &teeReadCloser{r, w}
}

type teeReadCloser struct {
	r io.ReadCloser
	w io.Writer
}

func (t *teeReadCloser) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		if n, err := t.w.Write(p[:n]); err != nil {
			return n, errors.WithStack(err)
		}
	}
	return
}

func (t *teeReadCloser) Close() error {
	return t.r.Close()
}
