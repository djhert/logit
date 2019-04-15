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
		return " -Error-"
	case 4:
		return " !Panic!"
	default:
		return ""
	}
}

const (
	MSG = iota
	WARN
	DEBUG
	ERROR
	PANIC
)

//Logger Logic
type Logger struct {
	TimeFormat string

	file io.WriteCloser

	log     chan msg
	closure chan bool
}

type msg struct {
	lvl Status
	s   string
}

func Start(f io.WriteCloser) (*Logger, error) {
	l := new(Logger)
	l.file = f
	l.log = make(chan msg)
	l.closure = make(chan bool)

	l.TimeFormat = "[2006/01/02 15:04:05.999999]"

	go logger(l)

	return l, nil
}

func logger(l *Logger) {
	defer l.file.Close()
	defer close(l.log)
	defer close(l.closure)

loop:
	for {
		select {
		case log := <-l.log:
			l.file.Write(genString(log, time.Now().Format(l.TimeFormat)))
		case <-l.closure:
			break loop
		}
	}
}

func genString(s msg, t string) []byte {
	return []byte(fmt.Sprintf("%s%s %s\n", t, s.lvl.String(), s.s))
}

func (l *Logger) Quit() {
	l.closure <- false
	<-l.closure
}

func (l *Logger) Log(e Status, a ...string) {
	l.log <- msg{lvl: e, s: strings.Join(a, " ")}
}

func (l *Logger) Logf(e Status, format string, a ...interface{}) {
	l.log <- msg{lvl: e, s: fmt.Sprintf(format, a...)}
}

func (l *Logger) LogError(o Status, e error) {
	l.Log(o, e.Error())
}
