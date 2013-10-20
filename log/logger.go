// Package log provides a logging facility.
package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

type level int

const (
	emergency level = iota
	alert
	crit
	err
	warn
	notice
	info
	debug
)

func (l level) String() string {
	return levels[l]
}

var levels = map[level]string{
	emergency: "emergency",
	alert:     "alert",
	crit:      "crit",
	err:       "err",
	warn:      "warn",
	notice:    "notice",
	info:      "info",
	debug:     "debug",
}

func NewLogger(levelStr string, timeFormat string, writer io.Writer) (Interface, error) {
	lvl, err := parseLevel(levelStr)
	if err != nil {
		return nil, err
	}
	return &logger{lvl, timeFormat, writer}, nil
}

func parseLevel(lvl string) (level, error) {
	for l, name := range levels {
		if name == lvl {
			return l, nil
		}
	}
	return 0, fmt.Errorf("unknown level: %s", lvl)
}

type logger struct {
	level      level
	timeFormat string
	writer     io.Writer
}

// Emergency logs the given arguments, and calls os.Exit(1) afterwards.
func (l *logger) Emergency(format string, args ...interface{}) {
	l.log(emergency, format, args...)
	os.Exit(1)
}

func (l *logger) Alert(format string, args ...interface{}) error {
	return l.logError(alert, format, args...)
}

func (l *logger) Crit(format string, args ...interface{}) error {
	return l.logError(crit, format, args...)
}

func (l *logger) Err(format string, args ...interface{}) error {
	return l.logError(err, format, args...)
}

func (l *logger) Warn(format string, args ...interface{}) {
	l.log(warn, format, args...)
}

func (l *logger) Notice(format string, args ...interface{}) {
	l.log(notice, format, args...)
}

func (l *logger) Info(format string, args ...interface{}) {
	l.log(info, format, args...)
}

func (l *logger) Debug(format string, args ...interface{}) {
	l.log(debug, format, args...)
}

func (l *logger) logError(lvl level, format string, args ...interface{}) error {
	l.log(lvl, format, args...)
	return fmt.Errorf(format, args...)
}

func (l *logger) log(lvl level, format string, args ...interface{}) {
	if lvl > l.level {
		return
	}

	t := time.Now()
	msg := fmt.Sprintf(format, args...)
	msg = fmt.Sprintf("%s [%s] %s\n", t.Format(l.timeFormat), levels[lvl], msg)
	if _, err := io.WriteString(l.writer, msg); err != nil {
		fmt.Printf("log error: %s: could not write to: %#v", err, l.writer)
	}
}
