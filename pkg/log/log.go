package log

import (
	"fmt"
	"os"
)

type Logger interface {
	// Debug prints the message to stderr only on debug mode
	Debug(format string, a ...any)

	// Info prints the message to stderr
	Info(format string, a ...any)
}

type logger struct {
	isDebug bool
}

func New(isDebug bool) Logger {
	return logger{
		isDebug: isDebug,
	}
}

// Debug prints the message to stderr only on debug mode
func (l logger) Debug(format string, a ...any) {
	if l.isDebug {
		userMsg := fmt.Sprintf(format, a...)
		fmt.Fprintf(os.Stderr, "[DEBUG] %s\n", userMsg)
	}
}

// Info prints the message to stderr
func (l logger) Info(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
}
