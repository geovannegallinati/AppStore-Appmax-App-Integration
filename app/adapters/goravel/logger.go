package goravel

import "github.com/goravel/framework/facades"

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (Logger) Warningf(format string, args ...any) {
	facades.Log().Warningf(format, args...)
}

func (Logger) Errorf(format string, args ...any) {
	facades.Log().Errorf(format, args...)
}
