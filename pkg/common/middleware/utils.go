package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"
	Bearer              = "Bearer"
	XTraceIdKey         = "x-trace-id"
	GTraceIdKey         = "trace-id"
	TraceIdKey          = "traceId"
	Stack               = "stack"
	App                 = "app"
	Service             = "service"
	Input               = "input"
	Body                = "body"
	Method              = "method"
	Path                = "path"
	Request             = "request"
	Response            = "response"
	Context             = "context"
	Error               = "error"
	UserReferenceIdKey  = "userReferenceId"
)

func ParseBearerToken(r *http.Request) (string, error) {
	tokenHeader := r.Header.Get(AuthorizationHeader)
	if tokenHeader == "" {
		return "", errors.New("authorization header is empty")
	}

	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		return "", fmt.Errorf("invalid authorization header: %s", tokenHeader)
	}

	tokenType := strings.TrimSpace(splitted[0])
	token := strings.TrimSpace(splitted[1])

	if tokenType != Bearer {
		return "", fmt.Errorf("incorrect token type '%s', expecting 'Bearer'", tokenType)
	}

	return token, nil
}

func GetUserReferenceId(ctx context.Context) (string, error) {
	val := ctx.Value(UserReferenceIdKey)
	if userReferenceId, ok := val.(string); ok {
		return userReferenceId, nil
	}
	return "", errors.New("context has not vote in it")
}
