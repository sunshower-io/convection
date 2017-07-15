package core

import "fmt"

var log = CreateLogger()

type Level int


func SetLogger(logger Logger) {
	log = logger
}



type Logger interface {
	SetLevel(l Level)

	Debugf(f string, vs ...interface{})

	Errorf(f string, vs ...interface{})
}

func CreateLogger() Logger {
	return new(DefaultLogger)
}

type DefaultLogger struct {
	Logger
}

func (l *DefaultLogger) Errorf(s string, args ...interface{}) {
	println(fmt.Sprintf(s, args))
}

func (l *DefaultLogger) Debugf(s string, args ...interface{}) {
	println(fmt.Sprintf(s, args))
}


type NoOpLogger struct {
	Logger
}

func (l *NoOpLogger) Errorf(s string, args ...interface{}) {
}

func (l *NoOpLogger) Debugf(s string, args ...interface{}) {
}
