package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/labstack/gommon/color"
)

type (
	Log struct {
		level  Level
		out    io.Writer
		err    io.Writer
		prefix string
		sync.Mutex
		color color.Color
		lock  bool
	}
	Level uint8
)

const (
	Trace = iota
	Debug
	Info
	Notice
	Warn
	Error
	Fatal
	Off = 10
)

var (
	levels = []string{
		color.Cyan("TRACE"),
		color.Blue("DEBUG"),
		color.Green("INFO"),
		color.Magenta("NOTICE"),
		color.Yellow("WARN"),
		color.Red("ERROR"),
		color.RedBg("FATAL"),
	}
)

func New(prefix string) (l *Log) {
	l = &Log{
		level:  Debug,
		out:    os.Stdout,
		err:    os.Stderr,
		prefix: prefix,
	}
	if runtime.GOOS == "windows" {
		l.lock = true
		color.Disable()
	}
	return
}

func (l *Log) SetLevel(v Level) {
	l.level = v
}

func (l *Log) SetOutput(w io.Writer) {
	l.out = w
	l.err = w

	switch w.(type) {
	case *os.File:
	case *bytes.Buffer:
		l.lock = true
		color.Disable()
	default:
		l.lock = false
	}
}

func (l *Log) Trace(i interface{}) {
	l.log(Trace, l.out, i)
}

func (l *Log) Debug(i interface{}) {
	l.log(Debug, l.out, i)
}

func (l *Log) Info(i interface{}) {
	l.log(Info, l.out, i)
}

func (l *Log) Notice(i interface{}) {
	l.log(Notice, l.out, i)
}

func (l *Log) Warn(i interface{}) {
	l.log(Warn, l.out, i)
}

func (l *Log) Error(i interface{}) {
	l.log(Error, l.err, i)
}

func (l *Log) Fatal(i interface{}) {
	l.log(Fatal, l.err, i)
}

func (l *Log) log(v Level, w io.Writer, i interface{}) {
	if l.lock {
		l.Lock()
		defer l.Unlock()
	}
	if v >= l.level {
		fmt.Fprintf(w, "%s|%s|%v\n", levels[v], l.prefix, i)
	}
}
