package app

import "log"

type logger struct {
}

func (l *logger) Log(v ...interface{}) {
	log.Println(v...)
}

func (l *logger) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l *logger) Error(v ...interface{}) {
	log.Println(v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
