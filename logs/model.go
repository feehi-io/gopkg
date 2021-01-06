package logs

import (
	"encoding/json"
	"time"
)

// NewCommonField create common field, often used for identify which machine generate this log row in cluster.
func NewCommonField(key string, value string) *commonField {
	return &commonField{
		Key:   key,
		Value: value,
	}
}

type commonField struct {
	Key   string `json:"ContextKey"`
	Value string `json:"value"`
}

func (field commonField) MarshalJSON() ([]byte, error) {
	m := map[string]string{field.Key: field.Value}
	return json.Marshal(m)
}

// MessageHeader collect log message infos.
type MessageHeader struct {
	Level   Severity  `json:"level"`
	TraceID string    `json:"trace_id"`
	Time    time.Time `json:"time"`
	Line    int       `json:"line"`
	File    string    `json:"file"`
}

// Content full log infos.
type Content struct {
	Headers MessageHeader `json:"headers"`
	Message string        `json:"message"`
	Fields  []Field       `json:"fields,omitempty"`
}
