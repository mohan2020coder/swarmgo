package logging

import (
    "log"
    "os"
)

type Logger struct {
    *log.Logger
}

func NewLogger(path string) *Logger {
    f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("open log: %v", err)
    }
    return &Logger{Logger: log.New(f, "", log.LstdFlags|log.Lmicroseconds)}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
    l.Logger.Fatalf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
    l.Logger.Printf("ERROR: "+format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
    l.Logger.Printf("INFO: "+format, v...)
}
