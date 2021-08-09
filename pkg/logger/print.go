package logger

import (
	"fmt"
	"log"
)

func Debug(template string, args ...interface{}) {
	content := MakeLine(template, args...)
	if Sugar != nil {
		Sugar.Debug(content)
	} else {
		log.Print(content)
	}
}

func Info(template string, args ...interface{}) {
	content := MakeLine(template, args...)
	if Sugar != nil {
		Sugar.Info(content)
	} else {
		log.Print(content)
	}
}

func Error(template string, args ...interface{}) {
	content := MakeLine(template, args...)
	if Sugar != nil {
		Sugar.Error(content)
	} else {
		log.Print(content)
	}
}

func MakeLine(template string, args ...interface{}) (content string) {
	content = fmt.Sprintf(template, args...)
	return
}
