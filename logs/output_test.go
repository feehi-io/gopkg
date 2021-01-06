package logs

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testLogFile = os.TempDir() + "/test.log"

type outputWriteError struct {
}

func (o *outputWriteError) Flush() error {
	_, err := o.Write([]byte(""))
	return err
}

func (o *outputWriteError) Write(p []byte) (int, error) {
	return 0, errors.New("mock write error return")
}

func (o *outputWriteError) IsLevelNeedRecord(s Severity) bool {
	return true
}

func TestOutput_Write(t *testing.T) {
	output := NewOutPut(AllSeverities, &bytes.Buffer{})
	n, err := output.Write([]byte("test"))
	assert.Nil(t, err)
	assert.Equal(t, 4, n)
}

func TestOutput_Flush(t *testing.T) {
	output := NewOutPut(AllSeverities, &bytes.Buffer{})
	err := output.Flush()
	assert.Nil(t, err)
}

func TestNewOutPut(t *testing.T) {
	ot := NewOutPut(AllSeverities, &bytes.Buffer{})
	o, ok := ot.(*output)
	assert.True(t, ok)
	assert.Equal(t, AllSeverities, o.Levels)
}

func TestNewStdOutOutput(t *testing.T) {
	ot := NewStdOutOutput(AllSeverities)
	o, ok := ot.(*output)
	assert.True(t, ok)
	assert.Equal(t, AllSeverities, o.Levels)
}

func TestNewFileOutput(t *testing.T) {
	testCases := []struct {
		Input struct {
			Severities []Severity
			LogFile    string
		}
		Expected struct {
			Severities []Severity
			Error      bool
		}
		Message string
	}{
		{
			Input: struct {
				Severities []Severity
				LogFile    string
			}{Severities: AllSeverities, LogFile: testLogFile},
			Expected: struct {
				Severities []Severity
				Error      bool
			}{Severities: AllSeverities, Error: false},
			Message: "log file open successfully",
		},
		{
			Input: struct {
				Severities []Severity
				LogFile    string
			}{Severities: AllSeverities, LogFile: "/////error.txt"},
			Expected: struct {
				Severities []Severity
				Error      bool
			}{Severities: nil, Error: true},
			Message: "log file open with error",
		},
	}
	for _, testCase := range testCases {
		ot, err := NewFileOutput(testCase.Input.Severities, testCase.Input.LogFile)
		if testCase.Expected.Error {
			assert.NotNil(t, err, testCase.Message)
			_, ok := ot.(*output)
			assert.False(t, ok)
		} else {
			assert.Nil(t, err, testCase.Message)
			o, ok := ot.(*output)
			assert.True(t, ok)
			assert.Equal(t, testCase.Expected.Severities, o.Levels, testCase.Message)
		}
	}
	err := os.Remove(testLogFile)
	if err != nil {
		fmt.Println("remove test log file failed, plesase remove it yourself:", testLogFile)
	}
}

func TestOutput_isLevelNeedRecord(t *testing.T) {
	output := NewOutPut([]Severity{DebugLog, InfoLog}, &bytes.Buffer{})
	testCases := []struct {
		Input    Severity
		Expected bool
	}{
		{
			Input:    DebugLog,
			Expected: true,
		},
		{
			Input:    InfoLog,
			Expected: true,
		},
		{
			Input:    WarningLog,
			Expected: false,
		},
		{
			Input:    ErrorLog,
			Expected: false,
		},
		{
			Input:    FatalLog,
			Expected: false,
		},
	}
	for _, testCase := range testCases {
		result := output.IsLevelNeedRecord(testCase.Input)
		assert.Equal(t, testCase.Expected, result, fmt.Sprintf("isLevelNeedRecord input %s expected %t but got %t", severityName[testCase.Input], testCase.Expected, result))
	}
}
