// Package httpserver implements HTTP server.
package httpserver

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/chindada/leopard/pkg/portscan"
)

const (
	_defaultHost              = ""
	_defaultPort              = "80"
	_defaultReadTimeout       = 5 * time.Second
	_defaultReadHeaderTimeout = 5 * time.Second
	_defaultWriteTimeout      = 5 * time.Minute
	_defaultShutdownTimeout   = 3 * time.Second
)

// Server -.
type Server struct {
	srv *http.Server

	host string
	port string

	logger   Logger
	keyPath  string
	certPath string
}

// New -.
func New(handler http.Handler, opts ...Option) *Server {
	s := &Server{
		srv: &http.Server{
			ErrorLog:          slog.NewLogLogger(slog.NewTextHandler(io.Discard, nil), slog.LevelInfo),
			Handler:           handler,
			ReadHeaderTimeout: _defaultReadHeaderTimeout,
			ReadTimeout:       _defaultReadTimeout,
			WriteTimeout:      _defaultWriteTimeout,
		},
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.host == "" {
		s.host = _defaultHost
	}
	if s.port == "" {
		s.port = _defaultPort
	}
	s.srv.Addr = net.JoinHostPort(s.host, s.port)
	return s
}

func (s *Server) GetListenPort() string {
	_, port, _ := net.SplitHostPort(s.srv.Addr)
	return port
}

func (s *Server) Start() error {
	if err := s.tryListen(); err != nil {
		return err
	}
	return nil
}

func (s *Server) StartWithRandomPort() error {
	scanner := portscan.NewPortScan()
	Port(scanner.GetRandomPort())(s)
	s.srv.Addr = net.JoinHostPort(s.host, s.port)
	if err := s.tryListen(); err != nil {
		return err
	}
	return nil
}

func (s *Server) tryListen() error {
	errChan := make(chan error)
	go func() {
		if s.certPath == "" || s.keyPath == "" {
			err := s.srv.ListenAndServe()
			if err != nil {
				errChan <- err
			}
			return
		}
		err := s.srv.ListenAndServeTLS(s.certPath, s.keyPath)
		if err != nil {
			errChan <- err
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case err := <-errChan:
			return err
		case <-ticker.C:
			_, port, err := net.SplitHostPort(s.srv.Addr)
			if err != nil {
				return err
			}
			if s.getPortIsUsed(port) {
				s.Infof("HTTP Serve On %v", s.srv.Addr)
				return nil
			}
		}
	}
}

func (s *Server) getPortIsUsed(port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), 500*time.Millisecond)
	if err != nil && conn != nil {
		return false
	}
	if conn != nil {
		defer func() {
			if err = conn.Close(); err != nil {
				return
			}
		}()
		return true
	}
	ln, err := net.Listen("tcp", net.JoinHostPort("", port))
	if err != nil {
		return true
	}

	if ln != nil {
		defer func() {
			if err = ln.Close(); err != nil {
				return
			}
		}()
	}
	return false
}

func (s *Server) Infof(format string, args ...interface{}) {
	if s.logger != nil {
		s.logger.Infof(strings.ReplaceAll(format, "\n", ""), args...)
	} else {
		fmt.Printf(format, args...)
	}
}
