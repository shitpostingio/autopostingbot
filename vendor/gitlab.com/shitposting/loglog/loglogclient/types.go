package loglogclient

// LogType is a type denoting the various types of log entries loglog supports
//go:generate stringer -type=LogType
type LogType int

const (
	// INFO is a standard informational message
	INFO LogType = iota

	// ERROR is an exceptional message, should be handled with care
	ERROR

	// WARNING is an exceptional message, but more similar to an INFO message than to an ERROR
	WARNING
)
