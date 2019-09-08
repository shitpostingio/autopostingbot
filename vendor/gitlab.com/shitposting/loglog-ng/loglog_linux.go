// +build linux

package loglog

import (
	"log"
	"log/syslog"

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
	l := log.New(nil, ms+": ", log.Ldate|log.Ltime)
	logwriter, err := syslog.New(level, ms)
	if err != nil {
		return nil, xerrors.Errorf("cannot create logger with level %s: %w", level, err)
	}

	l.SetOutput(logwriter)

	return l, nil
}
