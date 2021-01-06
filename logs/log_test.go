package logs

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var curFile = "log_test.go"

func TestInit(t *testing.T) {
	assert.NotNil(t, log)
}

func TestSetOutputs(t *testing.T) {
	stdOutputs := []Output{
		NewStdOutOutput(AllSeverities),
		NewStdOutOutput([]Severity{DebugLog}),
	}
	SetOutputs(stdOutputs...)

	assert.Equal(t, log.options.outputs, stdOutputs)
}

func TestSetCommonFields(t *testing.T) {
	commonField := NewCommonField("instance", "machineA")
	SetCommonFields(commonField)
	assert.Equal(t, log.options.commonFields[0], commonField)
}

func TestSetDirHeader(t *testing.T) {
	SetDirHeader(true)
	assert.Equal(t, log.options.addDirHeader, true)
	SetDirHeader(false)
	assert.Equal(t, log.options.addDirHeader, false)
}

func TestSync(t *testing.T) {
	testCases := []struct {
		Mock struct {
			Outputs []Output
		}
		Expected []error
	}{
		{
			Mock:     struct{ Outputs []Output }{Outputs: []Output{NewStdOutOutput(AllSeverities)}},
			Expected: nil,
		},
		{
			Mock:     struct{ Outputs []Output }{Outputs: []Output{NewStdOutOutput(AllSeverities), &outputWriteError{}}},
			Expected: []error{errors.New("mock write error return")},
		},
		{
			Mock:     struct{ Outputs []Output }{Outputs: []Output{&outputWriteError{}, &outputWriteError{}}},
			Expected: []error{errors.New("mock write error return"), errors.New("mock write error return")},
		},
	}
	for _, testCase := range testCases {
		log.options.outputs = testCase.Mock.Outputs
		errs := Sync()
		if testCase.Expected == nil {
			assert.Nil(t, errs)
		} else {
			assert.Equal(t, testCase.Expected, errs)
		}
		testInitLogging()
	}
}

func TestDebug(t *testing.T) {
	testInitLogging()
	Debug(context.Background(), "test Debug", String("category", "debug category"))
	Sync()
	assert.Contains(t, outputCollects.String(), curFile)
	assert.Contains(t, outputCollects.String(), "test Debug")
	assert.Contains(t, outputCollects.String(), "debug category")
}

func TestDebugDepth(t *testing.T) {
	testInitLogging()
	DebugDepth(context.Background(), 1, "test DebugDepth", String("category", "debug category"), String("category2", "debug category2"))
	Sync()
	assert.Contains(t, outputCollects.String(), curFile)
	assert.Contains(t, outputCollects.String(), "test DebugDepth")
	assert.Contains(t, outputCollects.String(), "debug category")
	assert.Contains(t, outputCollects.String(), "debug category2")
}

func TestInfo(t *testing.T) {
	testInitLogging()
	Info(context.Background(), "test Info", Any("val", map[string]string{"iammapkey": "i am map value"}))
	Sync()
	assert.Contains(t, outputCollects.String(), curFile)
	assert.Contains(t, outputCollects.String(), "test Info")
	assert.Contains(t, outputCollects.String(), "i am map value")
}

func TestInfoDepth(t *testing.T) {
	testInitLogging()
	InfoDepth(context.Background(), 1, "test InfoDepth")
	Sync()
	assert.Contains(t, outputCollects.String(), curFile)
	assert.Contains(t, outputCollects.String(), "test InfoDepth")
}

func TestWarning(t *testing.T) {
	testInitLogging()
	type args struct {
		Domain string
		Info   struct {
			Name string
		}
	}
	argsInfo := &args{
		Domain: "www.feehi.com",
		Info:   struct{ Name string }{Name: "feehi-name"},
	}
	Warning(context.Background(), "test Warning", Any("args", argsInfo))
	Sync()
	assert.Contains(t, outputCollects.String(), curFile)
	assert.Contains(t, outputCollects.String(), "test Warning")
	assert.Contains(t, outputCollects.String(), "feehi-name")
	assert.Contains(t, outputCollects.String(), "www.feehi.com")
}

func TestWarningDepth(t *testing.T) {
	testInitLogging()
	WarningDepth(context.Background(), 1, "test WarningDepth", Any("category", "test_category"))
	Sync()
	assert.Contains(t, outputCollects.String(), curFile)
	assert.Contains(t, outputCollects.String(), "test WarningDepth")
	assert.Contains(t, outputCollects.String(), "test_category")
}

func TestError(t *testing.T) {
	testInitLogging()
	Error(context.Background(), "test Error", Err(errors.New("i am error")))
	Sync()
	assert.Contains(t, outputCollects.String(), curFile)
	assert.Contains(t, outputCollects.String(), "test Error")
	assert.Contains(t, outputCollects.String(), "error")
	assert.Contains(t, outputCollects.String(), "i am error")
}

func TestErrorDepth(t *testing.T) {
	testInitLogging()
	ErrorDepth(context.Background(), 1, "test ErrorDepth", String("category", "test category"))
	Sync()
	assert.Contains(t, outputCollects.String(), curFile)
	assert.Contains(t, outputCollects.String(), "test ErrorDepth")
	assert.Contains(t, outputCollects.String(), "test category")
}

func TestFatal(t *testing.T) {
	testInitLogging()
	Fatal(context.Background(), "test Fatal")
	Sync()
	assert.Contains(t, outputCollects.String(), curFile)
	assert.Contains(t, outputCollects.String(), "test Fatal")
}

func TestFatalDepth(t *testing.T) {
	testInitLogging()
	FatalDepth(context.Background(), 1, "test FatalDepth")
	Sync()
	assert.Contains(t, outputCollects.String(), curFile)
	assert.Contains(t, outputCollects.String(), "test FatalDepth")
}

func testInitLogging() {
	log = NewLogging()
	outputCollects = bytes.Buffer{}
	SetOutputs(NewOutPut(AllSeverities, &outputCollects))
}
