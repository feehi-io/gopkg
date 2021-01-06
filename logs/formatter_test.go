package logs

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewStringFormatter(t *testing.T) {
	stringFormatter := mockStringFormatter()
	assert.Equal(t, defaultTimeHeaderFormat(), stringFormatter.TimeFormat)
	assert.Equal(t, DefaultStringFormatTemplate, stringFormatter.Template)
	assert.False(t, stringFormatter.ToUTCTime)
}

func TestStringFormatter_Format(t *testing.T) {
	testCases := []struct {
		Mock struct {
			StringFormatter *StringFormatter
		}
		Input struct {
			CommonField []*commonField
			Content     *Content
		}
		Expected string
	}{
		{
			Mock: struct {
				StringFormatter *StringFormatter
			}{
				StringFormatter: NewStringFormatter(DefaultStringFormatTemplate, defaultTimeHeaderFormat(), false),
			},
			Input: struct {
				CommonField []*commonField
				Content     *Content
			}{CommonField: nil, Content: &Content{
				Headers: MessageHeader{
					Level:   DebugLog,
					TraceID: "test_trace_id",
					Line:    101,
					File:    "test.go",
				},
				Message: "test message",
				Fields:  nil,
			}},
			Expected: fmt.Sprintf("[%s test_trace_id {TIME} test.go:101] test message\n", severityName[DebugLog]),
		},
		{
			Mock: struct {
				StringFormatter *StringFormatter
			}{
				StringFormatter: NewStringFormatter(DefaultStringFormatTemplate, defaultTimeHeaderFormat(), false),
			},
			Input: struct {
				CommonField []*commonField
				Content     *Content
			}{CommonField: []*commonField{
				{Key: "instance", Value: "testMachine"},
				{Key: "language", Value: "Go"},
			}, Content: &Content{
				Headers: MessageHeader{
					Level:   DebugLog,
					TraceID: "test_trace_id",
					Line:    101,
					File:    "test.go",
				},
				Message: "test message",
				Fields:  nil,
			}},
			Expected: fmt.Sprintf("[instance:testMachine language:Go  %s test_trace_id {TIME} test.go:101] test message\n", severityName[DebugLog]),
		},
		{
			Mock: struct {
				StringFormatter *StringFormatter
			}{
				StringFormatter: NewStringFormatter(DefaultStringFormatTemplate, defaultTimeHeaderFormat(), false),
			},
			Input: struct {
				CommonField []*commonField
				Content     *Content
			}{CommonField: []*commonField{
				{Key: "instance", Value: "testMachine"},
				{Key: "language", Value: "Go"},
			}, Content: &Content{
				Headers: MessageHeader{
					Level:   DebugLog,
					TraceID: "test_trace_id",
					Line:    101,
					File:    "test.go",
				},
				Message: "test message",
				Fields: func() []Field {
					fields := make([]Field, 2)
					fields[0] = &filed{
						key: "category",
						val: "Go",
					}
					fields[1] = &filed{
						key: "account_id",
						val: "123",
					}
					return fields
				}(),
			}},
			Expected: fmt.Sprintf("[instance:testMachine language:Go  %s test_trace_id {TIME} test.go:101] test message {category:Go,account_id:123}\n", severityName[DebugLog]),
		},
	}

	for _, testCase := range testCases {
		cur := time.Now()
		testCase.Input.Content.Headers.Time = cur
		testCase.Expected = strings.Replace(testCase.Expected, "{TIME}", cur.Format(defaultTimeHeaderFormat()), 1)
		message := testCase.Mock.StringFormatter.Format(testCase.Input.CommonField, testCase.Input.Content)
		assert.Equal(t, testCase.Expected, string(message))
	}
}

func TestNewJSONFormatter(t *testing.T) {
	JSONFormatter := mockJSONFormatter()
	assert.NotNil(t, JSONFormatter)
}

func TestJsonFormatter_Format(t *testing.T) {
	testCases := []struct {
		Mock struct {
			JSONMarshal func(interface{}) ([]byte, error)
		}
		Input    *Content
		Expected struct {
			Error   bool
			Message *Content
		}
	}{
		{
			Mock: struct {
				JSONMarshal func(interface{}) ([]byte, error)
			}{JSONMarshal: json.Marshal},
			Input: mockContent(),
			Expected: struct {
				Error   bool
				Message *Content
			}{Error: false, Message: mockContent()},
		},
		{
			Mock: struct {
				JSONMarshal func(interface{}) ([]byte, error)
			}{JSONMarshal: func(i interface{}) ([]byte, error) {
				return nil, errors.New("marshal occur error")
			}},
			Input: &Content{
				Headers: MessageHeader{
					Level: DebugLog,
				},
				Message: "test message",
				Fields:  nil,
			},
			Expected: struct {
				Error   bool
				Message *Content
			}{Error: true, Message: nil},
		},
	}
	JSONFormatter := mockJSONFormatter()
	for k, testCase := range testCases {
		if k != 1 {
			continue
		}
		jsonMarshal = testCase.Mock.JSONMarshal
		var bytes []byte
		errStr := testCaptureSTDOutput(func() {
			bytes = JSONFormatter.Format(nil, testCase.Input)
		})

		if testCase.Expected.Error {
			assert.Equal(t, "marshal log to json err marshal occur error\n", errStr)
		} else {
			assert.Empty(t, errStr)
			message := new(Content)
			err := json.Unmarshal(bytes, &message)
			assert.Nil(t, err)
			assert.Equal(t, testCase.Expected.Message, message)
		}
	}

}

func mockStringFormatter() *StringFormatter {
	return NewStringFormatter(DefaultStringFormatTemplate, defaultTimeHeaderFormat(), false)
}

func mockJSONFormatter() *JSONFormatter {
	return NewJSONFormatter()
}

func mockContent() *Content {
	logTime, err := time.Parse("2006-01-02 15:04:05", "2020-11-20 00:00:00")
	if err != nil {
		panic(err)
	}
	return &Content{
		Headers: MessageHeader{
			Level:   DebugLog,
			TraceID: "test_trace_id",
			Time:    logTime,
			Line:    101,
			File:    "test.go",
		},
		Message: "test message",
		Fields:  nil,
	}
}
