// +build linux

package loglog

import (
	"io/ioutil"
	"log"
	"log/syslog"
	"os"
	"strings"

	"golang.org/x/xerrors"
)

// Setup setups loglog for usage, by setting an appropriate messageSource and
// initializing all the needed loggers.
func Setup(messageSource string) error {
	if messageSource == "" {
		return xerrors.Errorf("messageSource must be filled")
	}

	ms = messageSource

	var err error

	logNotice, err = loggerForLevel(syslog.LOG_NOTICE)
	if err != nil {
		return xerrors.Errorf("cannot setup notice logger: %w", err)
	}

	logWarn, err = loggerForLevel(syslog.LOG_WARNING)
	if err != nil {
		return xerrors.Errorf("cannot setup notice logger: %w", err)
	}

	logErr, err = loggerForLevel(syslog.LOG_ERR)
	if err != nil {
		return xerrors.Errorf("cannot setup notice logger: %w", err)
	}

	PostLog = func(format ...interface{}) {
		log.Println(format...)
	}

	return nil
}

// loggerForLevel returns an instance of log.Logger which writes to syslog,
// with a given priority level.
func loggerForLevel(level syslog.Priority) (*log.Logger, error) {
	l := log.New(os.Stdout, ms+": ", log.Ldate|log.Ltime)
	if !runningInWSL() && !runningInDocker() {
		logwriter, err := syslog.New(level, ms)
		if err != nil {
			return nil, xerrors.Errorf("cannot create logger with level %s: %w", level, err)
		}
		l.SetOutput(logwriter)
	}

	return l, nil
}

// runningInWSL detects whether loglog-ng is running in WSL,
// where syslog doesn't work properly.
func runningInWSL() bool {
	return fileContains("/proc/version", "Microsoft@Microsoft.com")
}

// runningInDocker detects whether loglog-ng is running in Docker,
// where syslog doesn't behave correctly.
func runningInDocker() bool {
	return fileContains("/proc/1/cgroup", "docker")
}

// fileContains return true when path contains pattern.
func fileContains(path, pattern string) bool {
	pv, err := ioutil.ReadFile(path)
	if err != nil {
		return true
	}

	return strings.Contains(string(pv), pattern)
}
