package logs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithSkipHeaders(t *testing.T) {
	option := WithSkipHeaders(true)
	o := options{
		skipHeaders: false,
	}
	option(&o)

	assert.True(t, o.skipHeaders)
}

func TestWithAddDirHeader(t *testing.T) {
	option := WithAddDirHeader(true)
	o := options{
		addDirHeader: false,
	}
	option(&o)

	assert.True(t, o.addDirHeader)
}

func TestWithCommonFields(t *testing.T) {
	option := WithCommonField("instance", "test instance")
	o := options{}
	option(&o)
	assert.Equal(t, 1, len(o.commonFields))
	assert.Equal(t, "instance", o.commonFields[0].Key)
	assert.Equal(t, "test instance", o.commonFields[0].Value)
}

func TestWithOutputs(t *testing.T) {
	stdOutput := NewStdOutOutput(AllSeverities)
	option := WithOutput(stdOutput)
	o := options{
		outputs: []Output{NewOutPut(AllSeverities, &bytes.Buffer{})},
	}
	option(&o)

	assert.Equal(t, 2, len(o.outputs))
	assert.Equal(t, o.outputs[1], stdOutput)
}

func TestWithFormatter(t *testing.T) {
	option := WithFormatter(mockStringFormatter())
	o := options{
		formatter: NewJSONFormatter(),
	}
	option(&o)

	assert.Equal(t, mockStringFormatter(), o.formatter)
}

func TestWithMaxLogChanNum(t *testing.T) {
	option := WithMaxLogChanNum(1)
	o := options{
		maxLogChanNum: 10,
	}
	option(&o)

	assert.Equal(t, 1, o.maxLogChanNum)
}
