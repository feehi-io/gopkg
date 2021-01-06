package logs

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var curFileName string
var outputCollects bytes.Buffer

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("get file name failed")
	}
	curFileName = filepath.Base(file)
}

func TestNewLogging(t *testing.T) {
	l := NewLogging()
	if l == nil {
		t.Errorf("l is nil")
	}
}

func TestLoggingT_header(t *testing.T) {
	defer func() {
		runtimeCaller = runtime.Caller
	}()
	l := NewLogging()
	testCases := []struct {
		Mock struct {
			RuntimeCaller func(skip int) (pc uintptr, file string, line int, ok bool)
			TimeNow       func() time.Time
			AddDirHeader  bool
		}
		Input struct {
			Ctx context.Context
			S   Severity
		}
		Expected MessageHeader
	}{
		{
			Mock: struct {
				RuntimeCaller func(skip int) (pc uintptr, file string, line int, ok bool)
				TimeNow       func() time.Time
				AddDirHeader  bool
			}{RuntimeCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
				var a int
				return uintptr(a), "", 0, false
			}, TimeNow: nil, AddDirHeader: false},
			Input: struct {
				Ctx context.Context
				S   Severity
			}{Ctx: context.WithValue(context.Background(), TraceIDIdentifier, "test_trace_id_1"), S: DebugLog},
			Expected: MessageHeader{
				Level:   DebugLog,
				TraceID: "test_trace_id_1",
				Line:    1,
				File:    "???",
			},
		},
		{
			Mock: struct {
				RuntimeCaller func(skip int) (pc uintptr, file string, line int, ok bool)
				TimeNow       func() time.Time
				AddDirHeader  bool
			}{RuntimeCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
				var a int
				return uintptr(a), "/feehi/golang/src/test/a.go", 6, true
			}, TimeNow: nil, AddDirHeader: false},
			Input: struct {
				Ctx context.Context
				S   Severity
			}{Ctx: context.WithValue(context.Background(), TraceIDIdentifier, "test_trace_id_2"), S: DebugLog},
			Expected: MessageHeader{
				Level:   DebugLog,
				TraceID: "test_trace_id_2",
				Line:    6,
				File:    "a.go",
			},
		},
		{
			Mock: struct {
				RuntimeCaller func(skip int) (pc uintptr, file string, line int, ok bool)
				TimeNow       func() time.Time
				AddDirHeader  bool
			}{RuntimeCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
				var a int
				return uintptr(a), "/feehi/golang/src/test/a.go", 26, true
			}, TimeNow: nil, AddDirHeader: true},
			Input: struct {
				Ctx context.Context
				S   Severity
			}{Ctx: context.WithValue(context.Background(), TraceIDIdentifier, "test_trace_id_3"), S: DebugLog},
			Expected: MessageHeader{
				Level:   DebugLog,
				TraceID: "test_trace_id_3",
				Line:    26,
				File:    "test/a.go",
			},
		},
	}

	for _, testCase := range testCases {
		cur := time.Now()
		timeNow = func() time.Time {
			return cur
		}
		runtimeCaller = testCase.Mock.RuntimeCaller
		l.options.addDirHeader = testCase.Mock.AddDirHeader
		testCase.Expected.Time = cur
		messageHeader := l.header(testCase.Input.Ctx, testCase.Input.S, 0)
		assert.Equal(t, testCase.Expected, messageHeader)
	}
}

func TestLoggingT_Debug(t *testing.T) {
	l := testNewLogging()
	l.Debug(context.Background(), "test Debug", String("category", "debug category"))
	l.Sync()
	assert.Contains(t, outputCollects.String(), curFileName)
	assert.Contains(t, outputCollects.String(), "test Debug")
	assert.Contains(t, outputCollects.String(), "debug category")
}

func TestLoggingT_DebugDepth(t *testing.T) {
	l := testNewLogging()
	l.DebugDepth(context.Background(), 0, "test DebugDepth", String("category", "debug category"), String("category2", "debug category2"))
	l.Sync()
	assert.Contains(t, outputCollects.String(), curFileName)
	assert.Contains(t, outputCollects.String(), "test DebugDepth")
	assert.Contains(t, outputCollects.String(), "debug category")
	assert.Contains(t, outputCollects.String(), "debug category2")
}

func TestLoggingT_Info(t *testing.T) {
	l := testNewLogging()
	l.Info(context.Background(), "test Info", Any("val", map[string]string{"iammapkey": "i am map value"}))
	l.Sync()
	assert.Contains(t, outputCollects.String(), curFileName)
	assert.Contains(t, outputCollects.String(), "test Info")
	assert.Contains(t, outputCollects.String(), "i am map value")
}

func TestLoggingT_InfoDepth(t *testing.T) {
	l := testNewLogging()
	l.InfoDepth(context.Background(), 0, "test InfoDepth")
	l.Sync()
	assert.Contains(t, outputCollects.String(), curFileName)
	assert.Contains(t, outputCollects.String(), "test InfoDepth")
}

func TestLoggingT_Warning(t *testing.T) {
	l := testNewLogging()
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
	l.Warning(context.Background(), "test Warning", Any("args", argsInfo))
	l.Sync()
	assert.Contains(t, outputCollects.String(), curFileName)
	assert.Contains(t, outputCollects.String(), "test Warning")
	assert.Contains(t, outputCollects.String(), "feehi-name")
	assert.Contains(t, outputCollects.String(), "www.feehi.com")
}

func TestLoggingT_WarningDepth(t *testing.T) {
	l := testNewLogging()
	l.WarningDepth(context.Background(), 0, "test WarningDepth", Any("category", "test_category"))
	l.Sync()
	assert.Contains(t, outputCollects.String(), curFileName)
	assert.Contains(t, outputCollects.String(), "test WarningDepth")
	assert.Contains(t, outputCollects.String(), "test_category")
}

func TestLoggingT_Error(t *testing.T) {
	l := testNewLogging()
	l.Error(context.Background(), "test Error", Err(errors.New("i am error")))
	l.Sync()
	assert.Contains(t, outputCollects.String(), curFileName)
	assert.Contains(t, outputCollects.String(), "test Error")
	assert.Contains(t, outputCollects.String(), "error")
	assert.Contains(t, outputCollects.String(), "i am error")
}

func TestLoggingT_ErrorDepth(t *testing.T) {
	l := testNewLogging()
	l.ErrorDepth(context.Background(), 0, "test ErrorDepth", String("category", "test category"))
	l.Sync()
	assert.Contains(t, outputCollects.String(), curFileName)
	assert.Contains(t, outputCollects.String(), "test ErrorDepth")
	assert.Contains(t, outputCollects.String(), "test category")
}

func TestLoggingT_Fatal(t *testing.T) {
	l := testNewLogging()
	l.Fatal(context.Background(), "test Fatal")
	l.Sync()
	assert.Contains(t, outputCollects.String(), curFileName)
	assert.Contains(t, outputCollects.String(), "test Fatal")
}

func TestLoggingT_FatalDepth(t *testing.T) {
	l := testNewLogging()
	l.FatalDepth(context.Background(), 0, "test FatalDepth")
	l.Sync()
	assert.Contains(t, outputCollects.String(), curFileName)
	assert.Contains(t, outputCollects.String(), "test FatalDepth")
}

func TestLoggingT_print(t *testing.T) {
	l := testNewLogging()
	l.print(context.Background(), InfoLog, "test print", Any("category", "test_category"), String("domain", "feehi.com"))
	l.Sync()
	assert.Contains(t, outputCollects.String(), "testing.go")
	assert.Contains(t, outputCollects.String(), "test print")
	assert.Contains(t, outputCollects.String(), "test_category")
	assert.Contains(t, outputCollects.String(), "feehi.com")
}

func TestLoggingT_printDepth(t *testing.T) {
	l := testNewLogging()
	l.printDepth(context.Background(), InfoLog, -1, "test printDepth")
	l.Sync()
	assert.Contains(t, outputCollects.String(), curFileName)
	assert.Contains(t, outputCollects.String(), "test printDepth")
}

func TestLoggingT_output(t *testing.T) {
	outputCollects := bytes.Buffer{}
	l := NewLogging(WithOutput(NewOutPut([]Severity{DebugLog}, &outputCollects)))
	l.output(context.Background(), DebugLog, -2, "test output", Any("category", "test_category"), String("domain", "feehi.com"))
	l.output(context.Background(), InfoLog, -2, "test output no record")
	l.Sync()
	assert.Contains(t, outputCollects.String(), curFileName)
	assert.Contains(t, outputCollects.String(), "test output")
	assert.Contains(t, outputCollects.String(), "test_category")
	assert.Contains(t, outputCollects.String(), "feehi.com")
	assert.NotContains(t, outputCollects.String(), "test output no record")

}

func TestLoggingT_sync_success(t *testing.T) {
	outputCollectsInfo := bytes.Buffer{}
	outputCollectsDebug := bytes.Buffer{}
	l := NewLogging(
		WithOutput(NewOutPut([]Severity{InfoLog}, &outputCollectsInfo)),
		WithOutput(NewOutPut([]Severity{DebugLog}, &outputCollectsDebug)),
	)
	l.Debug(context.Background(), "test write debug")
	l.Info(context.Background(), "test write info")
	err := l.Sync()
	assert.Nil(t, err)
	assert.Contains(t, outputCollectsInfo.String(), "test write info")
	assert.NotContains(t, outputCollectsInfo.String(), "test write debug")

	assert.NotContains(t, outputCollectsDebug.String(), "test write info")
	assert.Contains(t, outputCollectsDebug.String(), "test write debug")
}

func TestLoggingT_sync_error(t *testing.T) {
	outputCollects = bytes.Buffer{}
	l := NewLogging(WithOutput(&outputWriteError{}))

	var errs []error
	stdOut := testCaptureSTDOutput(func() {
		l.Debug(context.Background(), "test write debug")
		l.Info(context.Background(), "test write info")
		errs = l.Sync()
	})

	assert.Equal(t, errors.New("mock write error return"), errs[0])
	assert.Equal(t, "write to log error mock write error return \nwrite to log error mock write error return \n", stdOut)
}

func TestLoggingT_write(t *testing.T) {
	testCases := []struct {
		Mock struct {
			Logging *logging
			After   func(l *logging)
		}
		Expected struct {
			STDOutput string
		}
	}{
		{
			Mock: struct {
				Logging *logging
				After   func(l *logging)
			}{Logging: NewLogging(), After: func(l *logging) {
				close(l.contentChan)
			}},
			Expected: struct{ STDOutput string }{STDOutput: "channel been closed unexpected \n"},
		},
	}

	for _, testCase := range testCases {
		testCase.Mock.After(testCase.Mock.Logging)
		stdOutput := testCaptureSTDOutput(func() {
			testCase.Mock.Logging.write()
		})
		assert.Equal(t, testCase.Expected.STDOutput, stdOutput)
	}
}

func TestLoggingT_writeLog(t *testing.T) {
	outputCollects1 := bytes.Buffer{}
	outputCollects2 := bytes.Buffer{}
	testCases := []struct {
		Logging  *logging
		Input    *Content
		Expected struct {
			LogString1 string
			LogString2 string
			STDOutput  string
		}
	}{
		{
			Logging: NewLogging(WithCommonField("HostName", "lf"), WithOutput(NewOutPut([]Severity{ErrorLog}, &outputCollects1)), WithOutput(NewOutPut([]Severity{InfoLog}, &outputCollects2))),
			Input: &Content{
				Headers: MessageHeader{
					Level:   ErrorLog,
					TraceID: "test_trace_id",
					Line:    20,
					File:    "test.go",
				},
				Message: "test message",
				Fields:  nil,
			},
			Expected: struct {
				LogString1 string
				LogString2 string
				STDOutput  string
			}{LogString1: "[HostName:lf  ERROR test_trace_id 0001-01-01 00:00:00 test.go:20] test message\n", LogString2: "", STDOutput: ""},
		},
		{
			Logging: NewLogging(WithCommonField("HostName", "lf"), WithOutput(NewOutPut([]Severity{ErrorLog}, &outputCollects1)), WithOutput(&outputWriteError{})),
			Input: &Content{
				Headers: MessageHeader{
					Level:   ErrorLog,
					TraceID: "test_trace_id",
					Line:    20,
					File:    "test.go",
				},
				Message: "test message write error",
				Fields:  nil,
			},
			Expected: struct {
				LogString1 string
				LogString2 string
				STDOutput  string
			}{LogString1: "[HostName:lf  ERROR test_trace_id 0001-01-01 00:00:00 test.go:20] test message write error\n", LogString2: "", STDOutput: "write to log error mock write error return \n"},
		},
	}

	for _, testCase := range testCases {
		outputCollects1 = bytes.Buffer{}
		outputCollects2 = bytes.Buffer{}
		stdOutput := testCaptureSTDOutput(func() {
			testCase.Logging.writeLog(testCase.Input)
			testCase.Logging.Sync()
		})
		assert.Equal(t, testCase.Expected.STDOutput, stdOutput)
		assert.Equal(t, testCase.Expected.LogString1, outputCollects1.String())
		assert.Equal(t, testCase.Expected.LogString2, outputCollects2.String())
		assert.Equal(t, testCase.Expected.STDOutput, stdOutput)
	}
}

func TestLoggingT_Sync(t *testing.T) {
	testCases := []struct {
		Input struct {
			Output          Output
			Message         string
			MessageCategory string
		}
		Expected struct {
			Errors   []error
			Contains []string
		}
	}{
		{
			Input: struct {
				Output          Output
				Message         string
				MessageCategory string
			}{
				Output:          NewOutPut(AllSeverities, &outputCollects),
				Message:         "test sync",
				MessageCategory: "test category",
			},
			Expected: struct {
				Errors   []error
				Contains []string
			}{Errors: nil, Contains: []string{"test sync", "test category"}},
		},
		{
			Input: struct {
				Output          Output
				Message         string
				MessageCategory string
			}{Output: &outputWriteError{}, Message: "", MessageCategory: ""},
			Expected: struct {
				Errors   []error
				Contains []string
			}{Errors: []error{errors.New("mock write error return")}, Contains: []string{}},
		},
	}
	for _, testCase := range testCases {
		l := NewLogging(WithOutput(testCase.Input.Output))
		l.output(context.Background(), InfoLog, -2, testCase.Input.Message, String("category", testCase.Input.MessageCategory))
		errs := l.Sync()
		if testCase.Expected.Errors == nil {
			assert.Contains(t, outputCollects.String(), curFileName)
			assert.Contains(t, outputCollects.String(), "test sync")
			assert.Contains(t, outputCollects.String(), "test category")
		} else {
			assert.Equal(t, testCase.Expected.Errors, errs)
		}
	}
}

func testNewLogging() *logging {
	outputCollects = bytes.Buffer{}
	return NewLogging(WithOutput(NewOutPut(AllSeverities, &outputCollects)))
}

func testCaptureSTDOutput(f func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	return string(out)
}
