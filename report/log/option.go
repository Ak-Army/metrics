package log

import "github.com/Ak-Army/xlog"

type Option struct {
	level   xlog.Level
	message string
}

func WithLevel(level xlog.Level) Options {
	return func(option *Option) {
		option.level = level
	}
}

func WithMessage(message string) Options {
	return func(option *Option) {
		option.message = message
	}
}
