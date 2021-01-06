package logs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Formatter format log to a row message
type Formatter interface {
	Format(commonFields []*commonField, row *Content) []byte
}

// DefaultStringFormatTemplate default StringFormatter template
const DefaultStringFormatTemplate = "[{COMMON_FIELDS} {LEVEL} {TRACE_ID} {TIME} {FILE}:{LINE}] {MESSAGE} {FIELDS}"

// NewStringFormatter create a StringFormatter
func NewStringFormatter(template string, timeFormat string, toUTCTime bool) *StringFormatter {
	return &StringFormatter{
		Template:   template,
		TimeFormat: timeFormat,
		ToUTCTime:  toUTCTime,
	}
}

//StringFormatter format log to a row string
type StringFormatter struct {
	Template   string
	TimeFormat string
	ToUTCTime  bool
}

//Format format log to a string
func (f StringFormatter) Format(commonFields []*commonField, message *Content) []byte {
	var s string
	if len(commonFields) > 0 {
		commonFieldsStr := ""
		for _, commonField := range commonFields {
			commonFieldsStr += commonField.Key + ":" + commonField.Value + " "
		}
		s = strings.Replace(f.Template, "{COMMON_FIELDS}", commonFieldsStr, -1)
	} else {
		s = strings.Replace(f.Template, "{COMMON_FIELDS} ", "", -1)
	}
	s = strings.Replace(s, "{LEVEL}", severityName[message.Headers.Level], -1)
	s = strings.Replace(s, "{TRACE_ID}", message.Headers.TraceID, -1)
	s = strings.Replace(s, "{TIME}", message.Headers.Time.Format(f.TimeFormat), -1)
	s = strings.Replace(s, "{LINE}", strconv.Itoa(message.Headers.Line), -1)
	s = strings.Replace(s, "{FILE}", message.Headers.File, -1)
	s = strings.Replace(s, "{MESSAGE}", message.Message, -1)

	fields := ""
	for k, field := range message.Fields {
		if k == 0 {
			fields = "{"
		}
		fields += field.Key() + ":" + field.Value() + ","
		if len(message.Fields)-1 == k {
			fields = strings.TrimRight(fields, ",")
			fields += "}"
		}

	}
	if len(message.Fields) > 0 {
		s = strings.Replace(s, "{FIELDS}", fields, -1)
	} else {
		s = strings.Replace(s, " {FIELDS}", fields, -1)
	}

	if s[len(s)-1] != '\n' {
		s += "\n"
	}
	return []byte(s)
}

// NewJSONFormatter create a JSONFormatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// JSONFormatter format log to JSON string
type JSONFormatter struct {
}

var jsonMarshal = json.Marshal

// Format format log to JSON row
func (f JSONFormatter) Format(commonFields []*commonField, content *Content) []byte {
	var record = struct {
		*Content
		CommonFields []*commonField `json:"common_fields"`
	}{
		Content:      content,
		CommonFields: commonFields,
	}
	s, err := jsonMarshal(record)
	if err != nil {
		fmt.Println("marshal log to json err", err)
	}
	return append(s, '\n')
}
