package logit

import (
	"io"
	"os"
)

func OpenFile(f string) (io.WriteCloser, error) {
	return os.OpenFile(f, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0640)
}

type stdlog struct {
	io.WriteCloser
}

func TermLog() io.WriteCloser {
	return &stdlog{}
}

func (s stdlog) Close() error {
	return nil
}

func (s stdlog) Write(a []byte) (int, error) {
	return io.WriteString(os.Stdout, string(a))
}
