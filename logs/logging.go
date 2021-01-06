package logs

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// NewLogging create log instance
func NewLogging(opts ...Option) *logging {
	options := defaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	if len(options.outputs) <= 0 {
		options.outputs = defaultLogOutputs()
	}
	if options.commonFields != nil && len(options.commonFields) == 0 {
		options.commonFields = defaultCommonFields()
	}
	l := &logging{
		options: options,
	}
	l.notifySyncChan = make(chan struct{}, 0)
	l.syncFinishChan = make(chan []error, 0)
	l.contentChan = make(chan *Content, l.options.maxLogChanNum)
	go l.write()
	return l
}

type options struct {
	TraceIDIdentifier ContextKey
	skipHeaders       bool
	addDirHeader      bool
	outputs           []Output
	formatter         Formatter
	commonFields      []*commonField
	maxLogChanNum     int
}

type logging struct {
	options        options
	contentChan    chan *Content
	notifySyncChan chan struct{}
	syncFinishChan chan []error
}

var runtimeCaller = runtime.Caller
var timeNow = time.Now

func (l *logging) header(ctx context.Context, s Severity, depth int) MessageHeader {
	_, file, line, ok := runtimeCaller(4 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		if slash := strings.LastIndex(file, "/"); slash >= 0 {
			path := file
			file = path[slash+1:]
			if l.options.addDirHeader {
				if dirsep := strings.LastIndex(path[:slash], "/"); dirsep >= 0 {
					file = path[dirsep+1:]
				}
			}
		}
	}
	traceID, _ := ctx.Value(l.options.TraceIDIdentifier).(string)
	return MessageHeader{
		Level:   s,
		TraceID: traceID,
		Time:    timeNow(),
		Line:    line,
		File:    file,
	}
}

func (l *logging) Debug(ctx context.Context, message string, fields ...Field) {
	l.print(ctx, DebugLog, message, fields...)
}

func (l *logging) DebugDepth(ctx context.Context, depth int, message string, fields ...Field) {
	l.printDepth(ctx, DebugLog, depth, message, fields...)
}

func (l *logging) Info(ctx context.Context, message string, fields ...Field) {
	l.print(ctx, InfoLog, message, fields...)
}

func (l *logging) InfoDepth(ctx context.Context, depth int, message string, fields ...Field) {
	l.printDepth(ctx, InfoLog, depth, message, fields...)
}

func (l *logging) Warning(ctx context.Context, message string, fields ...Field) {
	l.print(ctx, WarningLog, message, fields...)
}

func (l *logging) WarningDepth(ctx context.Context, depth int, message string, fields ...Field) {
	l.printDepth(ctx, WarningLog, depth, message, fields...)
}

func (l *logging) Error(ctx context.Context, message string, fields ...Field) {
	l.print(ctx, ErrorLog, message, fields...)
}

func (l *logging) ErrorDepth(ctx context.Context, depth int, message string, fields ...Field) {
	l.printDepth(ctx, ErrorLog, depth, message, fields...)
}

func (l *logging) Fatal(ctx context.Context, message string, fields ...Field) {
	l.print(ctx, FatalLog, message, fields...)
}

func (l *logging) FatalDepth(ctx context.Context, depth int, message string, fields ...Field) {
	l.printDepth(ctx, FatalLog, depth, message, fields...)
}

func (l *logging) print(ctx context.Context, s Severity, message string, fields ...Field) {
	l.printDepth(ctx, s, 1, message, fields...)
}

func (l *logging) printDepth(ctx context.Context, s Severity, depth int, message string, fields ...Field) {
	l.output(ctx, s, depth, message, fields...)
}

func (l *logging) output(ctx context.Context, s Severity, depth int, message string, fields ...Field) {
	content := &Content{
		Headers: l.header(ctx, s, depth),
		Message: message,
		Fields:  fields,
	}

	for _, output := range l.options.outputs { //exists one output this log level, should send to channel
		if output.IsLevelNeedRecord(s) {
			l.contentChan <- content
			break
		} else {
			continue
		}
	}
}

func (l *logging) write() {
	for {
		var content *Content
		var ok bool
		select {
		case <-l.notifySyncChan:
			ok = true
			for {
				if len(l.contentChan) <= 0 {
					break
				}
				content, ok = <-l.contentChan
				if !ok {
					break
				}
				l.writeLog(content)
			}
			var errs []error
			for _, output := range l.options.outputs {
				err := output.Flush()
				if err != nil {
					errs = append(errs, err)
				}
			}
			l.syncFinishChan <- errs
			if !ok {
				fmt.Printf("channel been closed unexpected \n")
				return
			}
		case content, ok = <-l.contentChan:
			if !ok {
				fmt.Printf("channel been closed unexpected \n")
				return
			}
			l.writeLog(content)
		}
	}
}

func (l *logging) writeLog(content *Content) {
	buf := l.options.formatter.Format(l.options.commonFields, content)
	for _, output := range l.options.outputs {
		if !output.IsLevelNeedRecord(content.Headers.Level) {
			continue
		}
		_, err := output.Write(buf)
		if err != nil {
			fmt.Printf("write to log error %s \n", err)
		}
	}
}

func (l *logging) Sync() []error {
	l.notifySyncChan <- struct{}{}
	return <-l.syncFinishChan
}
