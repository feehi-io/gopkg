package logs

// Option create NewLogging can pass option values.
type Option func(*options)

// WithSkipHeaders whether skipped log MessageHeader.
func WithSkipHeaders(skipHeaders bool) Option {
	return func(o *options) {
		o.skipHeaders = skipHeaders
	}
}

// WithAddDirHeader whether add dir in log message.(examples/a.go or a.go)
func WithAddDirHeader(addDirHeader bool) Option {
	return func(o *options) {
		o.addDirHeader = addDirHeader
	}
}

// WithCommonField set common field.
// Can called multi times, will set several common field.
// Often used for identify which machine generate that log row.
func WithCommonField(key string, value string) Option {
	return func(o *options) {
		o.commonFields = append(o.commonFields, &commonField{Key: key, Value: value})
	}
}

// WithOutput set output.
// Can called multi times, will set several outputs.
func WithOutput(output Output) Option {
	return func(o *options) {
		o.outputs = append(o.outputs, output)
	}
}

// WithFormatter set log formatter.
// Which will determine the log row format. Such as JSON or string etc.
func WithFormatter(formatter Formatter) Option {
	return func(o *options) {
		o.formatter = formatter
	}
}

// WithMaxLogChanNum set max buffered channel logs length.
// When buffered channel over this nums to be consumed, they will be blocked for call Debug(ctx, message) Info(ctx, message)...
func WithMaxLogChanNum(maxLogChanNum int) Option {
	return func(o *options) {
		o.maxLogChanNum = maxLogChanNum
	}
}
