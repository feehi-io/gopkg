package logs

//Severity log level
type Severity int32

const (
	// DebugLog debug log
	DebugLog Severity = iota
	// InfoLog info log
	InfoLog
	// WarningLog warning log
	WarningLog
	// ErrorLog error log
	ErrorLog
	// FatalLog fatal log
	FatalLog
)

// AllSeverities all supported log levels
var AllSeverities = []Severity{
	DebugLog,
	InfoLog,
	WarningLog,
	ErrorLog,
	FatalLog,
}

var severityName = []string{
	DebugLog:   "DEBUG",
	InfoLog:    "INFO",
	WarningLog: "WARNING",
	ErrorLog:   "ERROR",
	FatalLog:   "FATAL",
}
