package Logit

import (
	"fmt"
	"io"
	"time"
)

//Current Time Struct
type CurrentTime struct {
	Day   int
	Month int
	Year  int
	H     int
	M     int
	S     int
	Ns    int
}

func (c *CurrentTime) Format() string {
	return fmt.Sprintf("[%02d/%02d/%04d %02d:%02d:%02d.%04d] ", c.Month, c.Day, c.Year, c.H, c.M, c.S, c.Ns)
}

func (c *CurrentTime) Update() {
	now := time.Now() //Needed to make sure that the timestamp is in sync
	c.H, c.M, c.S = now.Clock()
	c.Ns = (now.Nanosecond() / 100000) //Trim the excess

	var s time.Month
	c.Year, s, c.Day = now.Date()
	c.Month = int(s)
}

//Logger Logic
type Logger struct {
	file   io.WriteCloser
	log    chan string
	unlock chan byte

	isOpen  bool
	closure chan bool
}

func StartLogger(f io.WriteCloser) (*Logger, error) {
	l := new(Logger)
	l.file = f
	l.isOpen = true
	l.log = make(chan string)
	l.closure = make(chan bool)
	l.unlock = make(chan byte)

	go logger(l)

	return l, nil
}

func logger(l *Logger) {

	defer l.file.Close()
	defer close(l.unlock)
	defer close(l.log)
	defer close(l.closure)

	tim := new(CurrentTime)
	tim.Update()

	for l.isOpen {
		select {
		case log := <-l.log:
			tim.Update()
			l.file.Write(genString(log, tim.Format()))
			l.unlock <- 0
		case close := <-l.closure:
			tim.Update()
			l.isOpen = close
			l.file.Write(genString("Closing", tim.Format()))
			l.closure <- false
		}
	}
}

func genString(s string, t string) []byte {
	return []byte(fmt.Sprintf("%s %s\n", t, s))
}

func (l *Logger) Quit() {
	l.closure <- false
	<-l.closure
}

func (l *Logger) Log(log string) {
	l.log <- log
	<-l.unlock
}

func (l *Logger) Logf(format string, a ...interface{}) {
	l.log <- fmt.Sprintf(format, a...)
	<-l.unlock
}

func genError(err error) error {
	if err == nil {
		return nil
	} else {
		return fmt.Errorf("Unable to start logger: %s", err)
	}
}
