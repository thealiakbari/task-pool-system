package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/thealiakbari/task-pool-system/cmd/executor/docs"
	"github.com/thealiakbari/task-pool-system/pkg/common/config"
	"github.com/thealiakbari/task-pool-system/pkg/common/ginh"
	"github.com/thealiakbari/task-pool-system/pkg/common/response"
)

type Handler interface {
	RegisterRoutes(c *gin.RouterGroup)
}

type Server struct {
	router *gin.Engine
	conf   *config.AppConfig
}

func NewServer(conf *config.AppConfig, handlers ...Handler) *Server {
	r := ginh.NewGinEngine(conf.Mode)

	server := &Server{
		router: r,
		conf:   conf,
	}

	server.registerRoutes(handlers...)
	return server
}

func (s *Server) registerRoutes(handlers ...Handler) {
	subRouter := s.router.Group("/api/v1")
	for _, handler := range handlers {
		handler.RegisterRoutes(subRouter)
	}
}

func (s *Server) Start() {
	err := s.router.Run(s.conf.Core.Http.Address)
	if err != nil {
		panic(err)
	}
}

func (s *Server) SwaggerApi() {
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	docs.SwaggerInfo.Title = "Task Pool Service"
	docs.SwaggerInfo.Description = "Todo App Service: This is a Task Pool App service."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"
	if s.conf.Mode == config.ModeLocal {
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%v", s.conf.Core.Http.Port)
		docs.SwaggerInfo.Schemes = []string{"http"}
	} else {
		docs.SwaggerInfo.Schemes = []string{"https"}
	}
}

func (s *Server) HealthCheck() {
	s.router.GET("/ping", func(ctx *gin.Context) {
		response.OKResponse(ctx, map[string]string{"message": "pong"})
	})
}

func (s *Server) Shutdown(ctx context.Context) error {
	srv := &http.Server{
		Addr:    s.conf.Core.Http.Address,
		Handler: s.router,
	}

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	log.Println("HTTP server shut down gracefully.")
	return nil
}
