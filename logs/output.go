package logs

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Output log output
type Output interface {
	io.Writer
	Flush() error
	IsLevelNeedRecord(s Severity) bool
}

type output struct {
	Levels []Severity
	Buffer *bufio.Writer
}

func (o output) Write(p []byte) (n int, err error) {
	return o.Buffer.Write(p)
}

func (o output) Flush() error {
	return o.Buffer.Flush()
}

func (o output) IsLevelNeedRecord(s Severity) bool {
	for _, l := range o.Levels {
		if l == s {
			return true
		}
	}
	return false
}

// NewOutPut create a log output
func NewOutPut(levels []Severity, writer io.Writer) Output {
	return &output{
		Buffer: bufio.NewWriter(writer),
		Levels: levels,
	}
}

// NewFileOutput create a file log output
func NewFileOutput(levels []Severity, filename string) (Output, error) {
	fl, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return output{}, fmt.Errorf("open log file error: %s", err)
	}
	return NewOutPut(levels, fl), nil
}

// NewStdOutOutput create a STD log output
func NewStdOutOutput(levels []Severity) Output {
	return NewOutPut(levels, os.Stdout)
}
