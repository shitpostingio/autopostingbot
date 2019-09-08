// +build linux windows

package loglog

import (
	"log"
)

var (
	ms        string
	logNotice *log.Logger
	logWarn   *log.Logger
	logErr    *log.Logger

	// PostLog gets executed after each function invocation.
	// By default, PostLog is a function which calls log.Println() on its
	// argument.
	PostLog func(...interface{})
)

// Err must be used to log errors.
func Err(format ...interface{}) {
	logErr.Println(format...)
	PostLog(format...)
}

// Info must be used to log general informations.
func Info(format ...interface{}) {
	logNotice.Println(format...)
	PostLog(format...)
}

// Warn must be used when something isn't exactly an error, but must be noticed
// by the user.
func Warn(format ...interface{}) {
	logWarn.Println(format...)
	PostLog(format...)
}
