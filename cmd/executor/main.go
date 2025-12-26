package main

import (
	"context"
	logger "log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/thealiakbari/task-pool-system/cmd"
	"golang.org/x/sync/errgroup"
)

func main() {
	conf := cmd.Setup()

	ctx, cancel := context.WithCancel(conf.Ctx)
	defer cancel()

	var healthy int32 = 1
	errGroup, ctx := errgroup.WithContext(ctx)

	// NOTE: Run the http and gRPC Server
	server := httpServer(conf)
	// Handle OS signals for graceful shutdown
	errGroup.Go(func() error {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-sigCh:
			logger.Printf("Received signal: %v, shutting down...", sig)
			atomic.StoreInt32(&healthy, 0)

			shutdownCtx, cancel := context.WithTimeout(conf.Ctx, 5*time.Second)
			defer cancel()
			if err := server.Shutdown(shutdownCtx); err != nil {
				return err
			}
			cancel()
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	// Wait for all goroutines to finish
	if err := errGroup.Wait(); err != nil && atomic.LoadInt32(&healthy) == 1 {
		logger.Fatalf("Error occurred: %v", err)
	} else {
		conf.Logger.Info(nil, "Shutdown complete.")
	}
}

// @termsOfService  http://swagger.io/terms/
// @contact.name   Task Pool Application
// @contact.url    https://swagger.io/support
// @BasePath  /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description "Type 'Bearer TOKEN' to correctly set the Authorization Bearer"
func httpServer(conf *cmd.SetupConfig) *Server {
	server := NewServer(
		conf.Conf,
		conf.HttpAdaptorStorage.TaskAdaptor,
	)

	server.HealthCheck()
	server.SwaggerApi()
	go server.Start()

	return server
}
