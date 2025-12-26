package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const stackLines = 6

type stackPtr []uintptr

func callers(skip int, keep ...int) stackPtr {
	keepStackLines := stackLines
	if len(keep) > 0 {
		keepStackLines = keep[0]
	}

	pcs := make([]uintptr, keepStackLines, keepStackLines)
	n := runtime.Callers(skip, pcs[:])
	var st stackPtr = pcs[:n]
	return st
}

var pwd string

func stacks(skip int, servicePackageName string) []string {
	if pwd == "" {
		pwd, _ = os.Getwd()
		pwd += "/"
	}

	pc := callers(skip)
	frames := runtime.CallersFrames(pc[:])

	res := make([]string, 0, stackLines)

	for {
		frame, more := frames.Next()
		file := strings.Replace(frame.File, pwd, "", 1)
		function := strings.Replace(frame.Function, servicePackageName, "", 1)

		res = append(res, fmt.Sprintf("%s:%d %s", file, frame.Line, function))
		if !more {
			break
		}
	}
	return res
}

func Stacks(skip int, keep ...int) []string {
	pc := callers(skip, keep...)
	frames := runtime.CallersFrames(pc[:])
	keepStackLines := stackLines
	if len(keep) > 0 {
		keepStackLines = keep[0]
	}
	res := make([]string, 0, keepStackLines)

	for {
		frame, more := frames.Next()
		res = append(res, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		if !more || len(res) >= keepStackLines {
			break
		}
	}
	return res
}
