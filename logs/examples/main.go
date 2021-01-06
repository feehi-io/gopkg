package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/feehi.io/gopkg/logs"
)

func main() {
	//directly use
	//directly()

	//custom use
	customString()
	//customJSON()
}

func directly() {
	defer logs.Sync()
	fileOutput, err := logs.NewFileOutput([]logs.Severity{logs.InfoLog}, "log.txt")
	if err != nil {
		panic(err)
	}
	logs.SetOutputs(fileOutput, logs.NewStdOutOutput([]logs.Severity{logs.DebugLog, logs.InfoLog, logs.ErrorLog}))
	//logs.SetDirHeader(true)
	//logs.SetCommonFields([]*logs.commonField{{Key: "instance", Value: "test_instance"}})
	//logs.SetCommonFields(nil)

	ctx := context.WithValue(context.Background(), "trace_id", "22222222222222")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		logs.Debug(ctx, "i am debug, and will output in std not for file", logs.String("category", "special category"))
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		logs.Info(ctx, "i am info also in std and file")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		logs.InfoDepth(ctx, 1, "i am infodepth")
		wg.Done()
	}()
	err = generateError()
	logs.Error(ctx, "get user info failed", logs.Err(err))

	wg.Wait()
}

func generateError() error {
	return fmt.Errorf("error get user from db: %w", deepError())
}

func deepError() error {
	return errors.New("not found id")
}

func customString() {
	ctx := context.WithValue(context.Background(), "trace_id", "aaaaaaaaa123")
	formatter := logs.WithFormatter(
		logs.NewStringFormatter(logs.DefaultStringFormatTemplate, time.RFC3339Nano, true), //format log to a string like `[{LEVEL} {TIME} {FILE}:{LINE}] {DATA}`
	)
	l := logs.NewLogging(formatter, logs.WithOutput(logs.NewStdOutOutput([]logs.Severity{logs.DebugLog, logs.ErrorLog})), logs.WithCommonField("instance", "test_instance"), logs.WithCommonField("abc", "aaa"))
	defer func() {
		errs := l.Sync()
		if errs != nil {
			panic(errs)
		}
	}()
	l.Debug(ctx, "custom debug", logs.String("i_am_key", "i am value"))
	l.Error(ctx, "abc", logs.Err(errors.New("occur error")))
}

func customJSON() {
	ctx := context.WithValue(context.Background(), "trace_id", "aaaaaaaaa123")
	formatter := logs.WithFormatter(
		logs.NewJSONFormatter(), //format log to json
	)
	fl, err := logs.NewFileOutput(logs.AllSeverities, "b.txt")
	if err != nil {
		panic(err)
	}
	l := logs.NewLogging(formatter, logs.WithOutput(fl), logs.WithCommonField("instance", "test_instance"))
	defer func() {
		errs := l.Sync()
		if errs != nil {
			panic(errs)
		}
	}()
	l.Debug(ctx, "custom debug", logs.String("category", "debug category"))
}
