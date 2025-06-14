package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-nunu/nunu-layout-mcp/pkg/log"
	"net/http"
	"time"
)

type Server struct {
	*gin.Engine
	httpSrv *http.Server
	host    string
	port    int
	logger  *log.Logger
}

type Option func(s *Server)

func NewServer(engine *gin.Engine, logger *log.Logger, opts ...Option) *Server {
	s := &Server{
		Engine: engine,
		logger: logger,
		host:   "0.0.0.0", // 默认 host
		port:   8080,      // 默认 port
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithServerHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

func WithServerPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	s.logger.Sugar().Infof("Starting server at %s", addr)

	s.httpSrv = &http.Server{
		Addr:         addr,
		Handler:      s,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// 启动 HTTP 服务
	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Sugar().Errorf("HTTP server error: %s", err)
		}
	}()

	// 等待 context 取消
	<-ctx.Done()
	s.logger.Sugar().Info("Context canceled, shutting down server...")
	return s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	s.logger.Sugar().Info("Shutting down server gracefully...")
	if err := s.httpSrv.Shutdown(shutdownCtx); err != nil {
		s.logger.Sugar().Errorf("Server forced to shutdown: %v", err)
		return err
	}

	s.logger.Sugar().Info("Server exited properly")
	return nil
}
