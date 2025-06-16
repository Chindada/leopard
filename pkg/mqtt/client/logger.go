package client

import (
	"fmt"
	"os"
	"strings"
)

type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

type logger struct {
	base Logger
}

func NewLogger(base Logger) Logger {
	if base == nil {
		return &logger{}
	}
	return &logger{base: base}
}

func (l *logger) Infof(format string, args ...any) {
	if l.base != nil {
		l.base.Infof(format, args...)
	} else {
		if !strings.Contains(format, "\n") {
			format = fmt.Sprintf("%s\n", format)
		}
		fmt.Printf(format, args...)
	}
}

func (l *logger) Warnf(format string, args ...any) {
	if l.base != nil {
		l.base.Warnf(format, args...)
	} else {
		if !strings.Contains(format, "\n") {
			format = fmt.Sprintf("%s\n", format)
		}
		fmt.Printf(format, args...)
	}
}

func (l *logger) Errorf(format string, args ...any) {
	if l.base != nil {
		l.base.Errorf(format, args...)
	} else {
		if !strings.Contains(format, "\n") {
			format = fmt.Sprintf("%s\n", format)
		}
		fmt.Printf(format, args...)
	}
}

func (l *logger) Fatalf(format string, args ...any) {
	if l.base != nil {
		l.base.Fatalf(format, args...)
	} else {
		if !strings.Contains(format, "\n") {
			format = fmt.Sprintf("%s\n", format)
		}
		fmt.Printf(format, args...)
		os.Exit(-1)
	}
}
