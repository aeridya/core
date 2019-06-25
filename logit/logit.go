package logit

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type Status int

func (s Status) String() string {
	switch s {
	case 1:
		return " (WARNING)"
	case 2:
		return " *DEBUG*"
	case 3:
		return " NOTICE:"
	case 4:
		return " -Error-"
	case 5:
		return " !Panic!"
	default:
		return ""
	}
}

const (
	MSG Status = iota
	WARN
	DEBUG
	NOTICE
	ERROR
	PANIC
)

var (
	TimeFormat string

	file io.WriteCloser

	log     chan msg
	closure chan bool
)

type msg struct {
	lvl Status
	s   string
}

func Start(f io.WriteCloser) error {
	file = f
	log = make(chan msg, 1)
	closure = make(chan bool, 1)

	TimeFormat = "[2006/01/02 15:04:05.999999]"

	go logger()

	return nil
}

func logger() {
	defer file.Close()
	defer close(log)
	defer close(closure)

loop:
	for {
		select {
		case log := <-log:
			file.Write(genString(log, time.Now().Format(TimeFormat)))
		case <-closure:
			break loop
		}
	}
}

func genString(s msg, t string) []byte {
	return []byte(fmt.Sprintf("%s%s %s\n", t, s.lvl.String(), s.s))
}

func Quit() {
	closure <- false
	<-closure
}

func Log(e Status, a ...string) {
	log <- msg{lvl: e, s: strings.Join(a, " ")}
}

func Logf(e Status, format string, a ...interface{}) {
	log <- msg{lvl: e, s: fmt.Sprintf(format, a...)}
}

func LogError(o Status, e error) {
	Log(o, e.Error())
}
