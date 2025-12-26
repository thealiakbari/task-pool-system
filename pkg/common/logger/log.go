package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/google/uuid"
	configx "github.com/thealiakbari/task-pool-system/pkg/common/config"
	"github.com/thealiakbari/task-pool-system/pkg/common/middleware"
	"github.com/thealiakbari/task-pool-system/pkg/common/utiles"
)

type Option func(*logger)

type logger struct {
	slogger            *slog.Logger
	skipStack          int
	appName            string
	servicePackageName string
	service            *string
}

func New(mode, serviceName string, servicePackageName string, opts ...Option) (Logger, error) {
	log, err := new(mode, serviceName, servicePackageName, opts...)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

func new(mode, serviceName string, servicePackageName string, opts ...Option) (logger, error) {
	logger := logger{
		skipStack:          5,
		servicePackageName: servicePackageName + "/",
		appName:            serviceName,
	}
	for _, opt := range opts {
		opt(&logger)
	}

	logger.slogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if mode != configx.ModeProd {
		// TODO: Somthing production related...
	}

	return logger, nil
}

func (l *logger) ForService(service interface{}) Logger {
	log := l.forService(service)
	return &log
}

func (l *logger) forService(service interface{}) logger {
	serviceName := getType(service)
	return logger{
		slogger:            l.slogger,
		servicePackageName: l.servicePackageName,
		service:            serviceName,
		appName:            l.appName,
		skipStack:          l.skipStack,
	}
}

func (l logger) Info(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, slog.LevelInfo, msg, fieldsToSlogAttrs(fields)...)
}

func (l logger) Error(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, slog.LevelError, msg, fieldsToSlogAttrs(fields)...)
}

func (l logger) MethodError(ctx context.Context, input interface{}, msg string, attrs ...Field) {
	l.log(ctx, slog.LevelError, msg, fieldsToSlogAttrs(append(attrs, Field{Key: "input", Value: input}))...)
}

func (l logger) Debug(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, slog.LevelDebug, msg, fieldsToSlogAttrs(fields)...)
}

func (l logger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, slog.LevelWarn, msg, fieldsToSlogAttrs(fields)...)
}

func (l logger) Infof(ctx context.Context, template string, args ...interface{}) {
	l.log(ctx, slog.LevelInfo, fmt.Sprintf(template, args...))
}

func (l logger) Errorf(ctx context.Context, template string, args ...interface{}) {
	l.log(ctx, slog.LevelError, fmt.Sprintf(template, args...))
}

func (l logger) MethodErrorf(ctx context.Context, input interface{}, template string, args ...interface{}) {
	l.log(ctx, slog.LevelError, fmt.Sprintf(template, args...), slog.Any(middleware.Input, input))
}

func (l logger) Debugf(ctx context.Context, template string, args ...interface{}) {
	l.log(ctx, slog.LevelDebug, fmt.Sprintf(template, args...))
}

func (l logger) Warnf(ctx context.Context, template string, args ...interface{}) {
	l.log(ctx, slog.LevelWarn, fmt.Sprintf(template, args...))
}

func (l logger) Panicf(ctx context.Context, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	l.log(ctx, slog.LevelError, msg)
	panic(msg)
}

func getType(service interface{}) *string {
	if service == nil {
		return nil
	}

	if t := reflect.TypeOf(service); t.Kind() == reflect.Ptr {
		return utiles.Ptr("*" + t.Elem().Name())
	} else {
		return utiles.Ptr(t.Name())
	}
}

func getTraceId(ctx context.Context) *string {
	if ctx == nil {
		return nil
	}

	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		traceIdUuid, ok := ctx.Value(middleware.TraceIdKey).(uuid.UUID)
		if !ok {
			return nil
		}
		traceId = traceIdUuid.String()
	}
	return &traceId
}

func (l logger) log(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	traceId := getTraceId(ctx)

	outAttrs := []slog.Attr{
		slog.String(middleware.App, l.appName),
	}

	if len(attrs) > 0 {
		outAttrs = append(outAttrs, attrs...)
	}

	if traceId != nil {
		outAttrs = append(outAttrs, slog.String(middleware.TraceIdKey, *traceId))
	}

	if l.service != nil {
		outAttrs = append(outAttrs, slog.String(middleware.Service, *l.service))
	}

	if level == slog.LevelDebug || level == slog.LevelWarn || level == slog.LevelError {
		stks := stacks(l.skipStack, l.servicePackageName)
		outAttrs = append(outAttrs, slog.Any(middleware.Stack, stks))
	}

	l.slogger.LogAttrs(ctx, level, msg, outAttrs...)
}

func (l logger) CloneAsInfra() InfraLogger {
	return &infraLogger{
		logger_impl: logger{
			slogger:            l.slogger,
			servicePackageName: l.servicePackageName,
			appName:            l.appName,
			skipStack:          6,
		},
	}
}
