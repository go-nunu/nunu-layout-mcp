package mcp

import (
	"context"
	"errors"
	"github.com/go-nunu/nunu-layout-mcp/pkg/log"
	"github.com/mark3labs/mcp-go/server"
	"net/http"
	"time"
)

type Server struct {
	*server.MCPServer
	stdio      bool
	stdioOpts  []server.StdioOption
	httpAddr   string
	sseAddr    string
	httpSrv    *server.StreamableHTTPServer
	sseSrv     *server.SSEServer
	logger     *log.Logger
	middleware http.Handler
}

type Option func(*Server)

func NewServer(logger *log.Logger, opts ...Option) *Server {
	s := &Server{
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
func WithStdioSrv(status bool, opts ...server.StdioOption) Option {
	return func(s *Server) {
		s.stdio = status
		s.stdioOpts = opts
	}
}
func WithMCPSrv(srv *server.MCPServer) Option {
	return func(s *Server) {
		s.MCPServer = srv
	}
}
func WithSSESrv(addr string, srv *server.SSEServer) Option {
	return func(s *Server) {
		s.sseSrv = srv
		s.sseAddr = addr
	}
}

func WithStreamableHTTPSrv(addr string, srv *server.StreamableHTTPServer) Option {
	return func(s *Server) {
		s.httpSrv = srv
		s.httpAddr = addr
	}
}

func (s *Server) Start(ctx context.Context) error {
	if s.MCPServer == nil {
		return errors.New("mcp server not initialized")
	}
	if s.stdio {
		go func() {
			s.logger.Sugar().Info("Starting STDIO server...")
			if err := server.ServeStdio(s.MCPServer, s.stdioOpts...); err != nil {
				s.logger.Sugar().Errorf("STDIO server error: %v", err)
			}
		}()
	}
	if s.sseSrv != nil {
		go func() {
			s.logger.Sugar().Infof("Starting SSE server on %s...", s.sseAddr)
			if err := s.sseSrv.Start(s.sseAddr); err != nil {
				s.logger.Sugar().Errorf("SSE server error: %v", err)
			}
		}()
	}
	if s.httpSrv != nil {
		s.logger.Sugar().Infof("Starting StreamableHTTP server on %s...", s.httpAddr)
		if err := s.httpSrv.Start(s.httpAddr); err != nil {
			s.logger.Sugar().Errorf("StreamableHTTP server error: %v", err)
		}
	}
	// 等待 context 取消
	<-ctx.Done()
	s.logger.Sugar().Info("Context canceled, shutting down server...")
	return s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	// 设置优雅关闭的超时时间
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	s.logger.Info("Shutting down server gracefully...")

	var shutdownErr error

	// 尝试关闭 SSE 服务
	if err := s.sseSrv.Shutdown(shutdownCtx); err != nil {
		s.logger.Sugar().Errorf("Failed to shutdown SSE server: %v", err)
		shutdownErr = err
	}

	// 尝试关闭 HTTP 服务
	if err := s.httpSrv.Shutdown(shutdownCtx); err != nil {
		s.logger.Sugar().Errorf("Failed to shutdown HTTP server: %v", err)
		if shutdownErr == nil {
			shutdownErr = err
		}
	}

	if shutdownErr == nil {
		s.logger.Info("Server exited properly")
	} else {
		s.logger.Warn("Server shutdown encountered issues")
	}

	return shutdownErr
}
