package logger

type infraLogger struct {
	logger_impl logger
}

func NewInfra(mode, serviceName string, servicePackageName string, opts ...Option) (InfraLogger, error) {
	log, err := new(mode, serviceName, servicePackageName, opts...)
	if err != nil {
		return nil, err
	}

	log.skipStack = 6

	return &infraLogger{
		logger_impl: log,
	}, nil
}

func (i *infraLogger) ForService(service interface{}) InfraLogger {
	log := i.logger_impl.forService(service)
	return &infraLogger{
		logger_impl: log,
	}
}

func (i infraLogger) Debug(msg string, attrs ...Field) {
	i.logger_impl.Debug(nil, msg, attrs...)
}

func (i infraLogger) Debugf(template string, args ...interface{}) {
	i.logger_impl.Debugf(nil, template, args...)
}

func (i infraLogger) Error(msg string, attrs ...Field) {
	i.logger_impl.Error(nil, msg, attrs...)
}

func (i infraLogger) Errorf(template string, args ...interface{}) {
	i.logger_impl.Errorf(nil, template, args...)
}

func (i infraLogger) Info(msg string, attrs ...Field) {
	i.logger_impl.Info(nil, msg, attrs...)
}

func (i infraLogger) Infof(template string, args ...interface{}) {
	i.logger_impl.Infof(nil, template, args...)
}

func (i infraLogger) Panicf(template string, args ...interface{}) {
	i.logger_impl.Panicf(nil, template, args...)
}

func (i infraLogger) Warn(msg string, attrs ...Field) {
	i.logger_impl.Warn(nil, msg, attrs...)
}

func (i infraLogger) Warnf(template string, args ...interface{}) {
	i.logger_impl.Warnf(nil, template, args...)
}

func (i infraLogger) CloneAsLogger() Logger {
	return &logger{
		slogger:            i.logger_impl.slogger,
		servicePackageName: i.logger_impl.servicePackageName,
		appName:            i.logger_impl.appName,
		skipStack:          5,
	}
}
