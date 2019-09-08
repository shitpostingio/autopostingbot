// +build windows

package loglog

import (
	"log"
	"os"

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

	logNotice, err = loggerForLevel()
	if err != nil {
		return xerrors.Errorf("cannot setup notice logger: %w", err)
	}

	logWarn, err = loggerForLevel()
	if err != nil {
		return xerrors.Errorf("cannot setup notice logger: %w", err)
	}

	logErr, err = loggerForLevel()
	if err != nil {
		return xerrors.Errorf("cannot setup notice logger: %w", err)
	}

	PostLog = func(format ...interface{}) {
		// nothing for you, windows.
	}

	return nil
}

// loggerForLevel returns an instance of log.Logger which writes to syslog,
// with a given priority level.
func loggerForLevel() (*log.Logger, error) {
	l := log.New(os.Stdout, ms+": ", log.Ldate|log.Ltime)

	return l, nil
}
