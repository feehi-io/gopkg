package logs

import (
	"context"
)

var log *logging

func init() {
	log = NewLogging()
}

// SetOutputs set log outputs.
// Default is std output(os.Stdout).
// If you does not want to output log, can set it to nil.
func SetOutputs(outputs ...Output) {
	log.options.outputs = outputs
}

// SetCommonFields set global message fields.
// Default is HostName(os.Hostname()). If you want no common fields, just set it to nil.
// Often used for Cluster to identify which machine generate that log.
func SetCommonFields(commonFields ...*commonField) {
	log.options.commonFields = commonFields
}

// SetDirHeader set if need print file log directory in log message. Such as example/a.go or only print a.go.
// Default is false.
func SetDirHeader(dir bool) {
	log.options.addDirHeader = dir
}

// Sync sync log to outputs.
// Because we use buffered writer, so only buffer exceed a value they will really write to storage.
// When Sync called, will trigger all outputs write buffer to storage.
func Sync() []error {
	return log.Sync()
}

// Debug record debug log
func Debug(ctx context.Context, message string, fields ...Field) {
	log.DebugDepth(ctx, 1, message, fields...)
}

// DebugDepth record debug log with assigned code file depth
func DebugDepth(ctx context.Context, depth int, message string, fields ...Field) {
	log.DebugDepth(ctx, depth, message, fields...)
}

// Info record info log
func Info(ctx context.Context, message string, fields ...Field) {
	log.InfoDepth(ctx, 1, message, fields...)
}

// InfoDepth record info log with assigned code file depth
func InfoDepth(ctx context.Context, depth int, message string, fields ...Field) {
	log.InfoDepth(ctx, depth, message, fields...)
}

// Warning record warning log
func Warning(ctx context.Context, message string, fields ...Field) {
	log.WarningDepth(ctx, 1, message, fields...)
}

// WarningDepth record warning log with assigned code file depth
func WarningDepth(ctx context.Context, depth int, message string, fields ...Field) {
	log.WarningDepth(ctx, depth, message, fields...)
}

// Error record error log
func Error(ctx context.Context, message string, fields ...Field) {
	log.ErrorDepth(ctx, 1, message, fields...)
}

// ErrorDepth record error log with assigned code file depth
func ErrorDepth(ctx context.Context, depth int, message string, fields ...Field) {
	log.ErrorDepth(ctx, depth, message, fields...)
}

// Fatal record fatal log
func Fatal(ctx context.Context, message string, fields ...Field) {
	log.FatalDepth(ctx, 1, message, fields...)
}

// FatalDepth record fatal log with assigned code file depth
func FatalDepth(ctx context.Context, depth int, message string, fields ...Field) {
	log.FatalDepth(ctx, depth, message, fields...)
}
