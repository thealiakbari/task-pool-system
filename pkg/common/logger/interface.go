package logger

import (
	"context"
)

type Logger interface {
	ForService(service interface{}) Logger
	Info(ctx context.Context, msg string, attrs ...Field)
	Error(ctx context.Context, msg string, attrs ...Field)
	MethodError(ctx context.Context, input interface{}, msg string, attrs ...Field)
	Debug(ctx context.Context, msg string, attrs ...Field)
	Warn(ctx context.Context, msg string, attrs ...Field)
	Infof(ctx context.Context, template string, args ...interface{})
	Errorf(ctx context.Context, template string, args ...interface{})
	MethodErrorf(ctx context.Context, input interface{}, template string, args ...interface{})
	Debugf(ctx context.Context, template string, args ...interface{})
	Warnf(ctx context.Context, template string, args ...interface{})
	Panicf(ctx context.Context, template string, args ...interface{})

	CloneAsInfra() InfraLogger
}

type InfraLogger interface {
	ForService(service interface{}) InfraLogger
	Info(msg string, attrs ...Field)
	Error(msg string, attrs ...Field)
	Debug(msg string, attrs ...Field)
	Warn(msg string, attrs ...Field)
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Panicf(template string, args ...interface{})

	CloneAsLogger() Logger
}
