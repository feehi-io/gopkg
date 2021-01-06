package logs

import (
	"encoding/json"
	"fmt"
)

// Field is additional info helps log more info
type Field interface {
	Key() string
	Value() string
}

type filed struct {
	key string
	val interface{}
}

func (f filed) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		f.Key(): f.Value(),
	})
}

func (f *filed) Key() string {
	return f.key
}
func (f *filed) Value() string {
	return fmt.Sprintf("%s", f.val)
}

// String log additional string value
func String(key string, value string) Field {
	return &filed{
		key: key,
		val: value,
	}
}

// Any log any data type additional value
func Any(key string, value interface{}) Field {
	return &filed{
		key: key,
		val: fmt.Sprintf("%v", value),
	}
}

// Err log error additional value
func Err(err error) Field {
	return &filed{
		key: "error",
		val: err,
	}
}
