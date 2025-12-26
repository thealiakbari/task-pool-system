package ginh

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thealiakbari/task-pool-system/pkg/common/logger"
	"github.com/thealiakbari/task-pool-system/pkg/common/middleware"
	"github.com/thealiakbari/task-pool-system/pkg/common/response"
	"github.com/thealiakbari/task-pool-system/pkg/common/utiles"
	slog "log/slog"
)

var slogger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

func traceIdMiddlewarefunc(c *gin.Context) {
	traceIdStr := c.Request.Header.Get("X-Trace-Id")
	traceId, err := uuid.Parse(traceIdStr)
	if err != nil {
		traceId = uuid.New()
	}

	c.Set(middleware.XTraceIdKey, traceId)
	ctx := context.WithValue(c.Request.Context(), middleware.TraceIdKey, traceId)
	c.Request = c.Request.WithContext(ctx)
	c.Header(middleware.XTraceIdKey, traceId.String())
	c.Next()
}

func responseLoggerMiddleware(c *gin.Context) {
	var body any
	if c.Request.Body != nil {
		// Read the request body and reset it to its original state.
		requestBody, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		var bodyMap map[string]any
		err := json.Unmarshal(requestBody, &bodyMap)
		if err != nil {
			body = string(requestBody)
		} else {
			body = bodyMap
		}
	} else {
		body = nil
	}

	// Capture the original response writer.
	originalWriter := c.Writer
	// Create a custom writer to capture the response body.
	bodyCapture := &responseBodyCapture{ResponseWriter: originalWriter, body: bytes.NewBufferString("")}
	c.Writer = bodyCapture

	logResult := func() {
		var response any
		var resultMap map[string]any
		context := make(map[string]any)
		context["method"] = c.Request.Method
		context["path"] = c.Request.URL.Path
		context["params"] = utiles.SimplifyMap(c.Request.URL.Query())
		context["req-headers"] = utiles.SimplifyMap(c.Request.Header)
		context["res-headers"] = utiles.SimplifyMap(c.Writer.Header())
		context["status"] = c.Writer.Status()

		err := json.Unmarshal(bodyCapture.body.Bytes(), &resultMap)
		if err != nil {
			response = bodyCapture.body.String()
		} else {
			response = resultMap
		}

		status := "request "
		level := slog.LevelInfo
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			status += "success - "
		} else if c.Writer.Status() >= 400 && c.Writer.Status() < 500 {
			status += "client error - "
			level = slog.LevelError
		} else if c.Writer.Status() >= 500 {
			status += "server error - "
			level = slog.LevelError
		}

		traceId, _ := c.Get(middleware.XTraceIdKey)
		status += strconv.FormatInt(int64(c.Writer.Status()), 10)

		// Log the response information after the request is processed.
		slogger.LogAttrs(nil, level, status,
			slog.Any(middleware.TraceIdKey, traceId),
			slog.Any(middleware.Body, body),
			slog.Any(middleware.Response, response),
			slog.Any(middleware.Context, context),
		)
	}

	panicked := true
	defer func() {
		if r := recover(); r != nil || panicked {
			traceId, _ := c.Get(middleware.XTraceIdKey)

			slogger.Debug(
				"PANIC ",
				slog.Any(middleware.TraceIdKey, traceId),
				slog.String(middleware.Error, fmt.Sprintf("%v", r)),
				// The skip with 8 frames, come from the recovery functions from APM, Gin and Go
				slog.Any(middleware.Stack, logger.Stacks(8)),
			)
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				response.BaseResponse{
					Payload: nil,
					Meta: response.ErrResponse{
						Message: "Internal Server Error",
						Causes:  nil,
						Code:    500,
					},
				},
			)
			logResult()
		} else {
			logResult()
		}
	}()
	c.Next()
	panicked = false

	// Restore the original writer.
	c.Writer = originalWriter
}

// responseBodyCapture is a custom ResponseWriter that captures the response body.
type responseBodyCapture struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body and writes to the original writer.
func (w *responseBodyCapture) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func NewGinEngine(mode string) *gin.Engine {
	// Usingh New to drop the gin.Logger
	r := gin.New()
	// r.Use(recovery)
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}))

	r.Use(traceIdMiddlewarefunc)
	r.Use(responseLoggerMiddleware)

	return r
}
