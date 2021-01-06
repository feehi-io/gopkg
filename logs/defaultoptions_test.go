package logs

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultOptions(t *testing.T) {
	options := defaultOptions()
	assert.Equal(t, TraceIDIdentifier, options.TraceIDIdentifier)
	assert.Equal(t, defaultFormatter(), options.formatter)
	assert.Equal(t, false, options.addDirHeader)
}

func TestDefaultLogOutputs(t *testing.T) {
	outputs := defaultLogOutputs()
	assert.Equal(t, 1, len(outputs))
	assert.Equal(t, NewStdOutOutput(AllSeverities), outputs[0])
}

func TestDefaultTimeHeaderFormat(t *testing.T) {
	timeFormat := defaultTimeHeaderFormat()
	_, err := time.Parse(timeFormat, "2020-11-20 00:00:00")
	assert.Equal(t, nil, err)
}

func TestDefaultFormatter(t *testing.T) {
	defaultFormatter := defaultFormatter()
	_, ok := defaultFormatter.(*StringFormatter)
	assert.Equal(t, true, ok)
}

func TestDefaultCommonFields(t *testing.T) {
	testCases := []struct {
		MockOsHostname func() (string, error)
		Expected       []*commonField
	}{
		{
			MockOsHostname: func() (string, error) {
				return "", errors.New("error occur")
			},
			Expected: []*commonField{},
		},
		{
			MockOsHostname: func() (string, error) {
				return "test-host", nil
			},
			Expected: []*commonField{{Key: "HostName", Value: "test-host"}},
		},
	}
	for _, testCase := range testCases {
		osHostname = testCase.MockOsHostname
		defaultCommonFields := defaultCommonFields()
		assert.Equal(t, testCase.Expected, defaultCommonFields)
	}
}
