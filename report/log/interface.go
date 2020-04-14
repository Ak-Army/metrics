package log

type Options func(*Option)

type logReport interface {
	getValue() string
	getName() string
}
