package db

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/google/uuid"
	"github.com/thealiakbari/task-pool-system/pkg/common/logger"
	"github.com/thealiakbari/task-pool-system/pkg/common/middleware"
	"github.com/thealiakbari/task-pool-system/pkg/common/utiles"
	glog "gorm.io/gorm/logger"
)

const (
	highlightSql = true
	reset        = "\033[0m"
	red          = "\033[31m"
	magenta      = "\033[35m"
	green        = "\033[32m"
	blueBold     = "\033[34;1m"
	infoStr      = green + reset + green + "[info] " + reset
	warnStr      = blueBold + reset + magenta + "[warn] " + reset
	errStr       = magenta + reset + red + "[error] " + reset
)

type gormLogger struct {
	traceStacks bool
}

var _ glog.Interface = gormLogger{}

func newGormLogger(traceStacks bool) glog.Interface {
	return gormLogger{traceStacks: traceStacks}
}

func nowStr() string {
	now := time.Now()
	return fmt.Sprintf("%02d/%02d/%02d %02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
	)
}

func log(sql string, rows int64, level string, beginTime time.Time, elapsed time.Duration, traceId *string, traceStacks bool, err error) {
	x := bytes.NewBufferString(
		fmt.Sprintf(
			"[%s] %s elapsed: %v -- rows: %d",
			nowStr(),
			level,
			elapsed.String(),
			rows,
		),
	)
	if traceId != nil {
		x.WriteString(" -- traceId: ")
		x.WriteString(*traceId)
	}

	x.WriteString("\n")

	if highlightSql {
		quick.Highlight(x, sql, "postgresql", "terminal16m", "monokai")
	} else {
		x.WriteString(sql)
	}

	if traceStacks || err != nil {
		stack := logger.Stacks(7, 4)
		x.WriteString("\n")
		if err != nil {
			x.WriteString("Error: ")
			x.WriteString(err.Error())
			x.WriteString("\n")
		}
		for _, st := range stack {
			x.WriteString(st)
			x.WriteString("\n")
		}
	}

	println(x.String())
}

func getStack() string {
	stack := logger.Stacks(7, 4)
	return strings.Join(stack, "\n")
}

func (g gormLogger) LogMode(glog.LogLevel) glog.Interface {
	return g
}

func (g gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	// slowThreshold from gorm source code, in default logger this is a constant initiation
	slowThreshold := 200 * time.Millisecond
	sql, rows := fc()

	var traceId *string
	traceIdUuid, ok := ctx.Value(middleware.TraceIdKey).(uuid.UUID)
	if ok {
		traceId = utiles.Ptr(traceIdUuid.String())
	}

	level := infoStr
	if err != nil {
		level = errStr
	} else if elapsed > slowThreshold {
		level = warnStr
	}

	log(sql, rows, level, begin, elapsed, traceId, g.traceStacks, err)
}

func (g gormLogger) Warn(ctx context.Context, template string, args ...interface{}) {
	fmt.Printf(warnStr+" "+template+"\n%s", args, getStack())
}

func (g gormLogger) Error(ctx context.Context, template string, args ...interface{}) {
	fmt.Printf(errStr+" "+template+"\n%s", args, getStack())
}

func (g gormLogger) Info(ctx context.Context, template string, args ...interface{}) {
	fmt.Printf(infoStr+" "+template+"\n%s", args, getStack())
}
