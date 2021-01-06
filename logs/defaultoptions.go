package logs

import (
	"os"
)

func defaultOptions() options {
	return options{
		TraceIDIdentifier: TraceIDIdentifier,
		formatter:         defaultFormatter(),
		commonFields:      []*commonField{},
		addDirHeader:      false,
		maxLogChanNum:     1000,
	}
}

func defaultLogOutputs() []Output {
	outputs := []Output{NewStdOutOutput(AllSeverities)}

	return outputs
}

func defaultTimeHeaderFormat() string {
	return "2006-01-02 15:04:05"
}

func defaultFormatter() Formatter {
	return NewStringFormatter(DefaultStringFormatTemplate, defaultTimeHeaderFormat(), false)
}

var osHostname = os.Hostname

func defaultCommonFields() []*commonField {
	commonFields := make([]*commonField, 0)
	hostName, err := osHostname()
	if err == nil {
		commonFields = append(commonFields, &commonField{
			Key:   "HostName",
			Value: hostName,
		})
	}
	return commonFields
}
