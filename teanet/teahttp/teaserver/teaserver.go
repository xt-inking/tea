package teaserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xt-inking/tea/tealog"
)

type Server struct {
	router     router
	logger     Logger
	config     *Config
	onShutdown []func()
}

type Logger interface {
	Close()
	Debug(ctx context.Context, msg string)
	Info(ctx context.Context, msg string)
	Error(ctx context.Context, msg string)
}

type Config struct {
	Addr                string
	ReadTimeout         time.Duration
	ReadHeaderTimeout   time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
	MaxHeaderBytes      int
	MaxHeaderValueCount int
	ShutdownTimeout     time.Duration
}

func New(opts ...option) *Server {
	s := &Server{
		router: newRouter(),
		logger: tealog.New(
			tealog.NewRecordText,
			tealog.WriterCloserFile[tealog.RecordText](tealog.FileDir, "http"),
		),
		config: &Config{
			Addr:                ":8080",
			ReadTimeout:         time.Second * 15,
			ReadHeaderTimeout:   time.Second * 5,
			WriteTimeout:        time.Second * 30,
			IdleTimeout:         time.Second * 120,
			MaxHeaderBytes:      http.DefaultMaxHeaderBytes,
			MaxHeaderValueCount: http.DefaultMaxHeaderValueCount,
			ShutdownTimeout:     time.Second * 60,
		},
		onShutdown: nil,
	}
	for _, o := range opts {
		o.apply(s)
	}
	return s
}

func (s *Server) Router() *router {
	return &s.router
}

func (s *Server) RegisterOnShutdown(f func()) {
	s.onShutdown = append(s.onShutdown, f)
}

func (s *Server) Serve() {
	defer func() {
		s.logger.Close()
		for _, f := range s.onShutdown {
			f()
		}
	}()
	server := &http.Server{
		Addr:                s.config.Addr,
		Handler:             s,
		ReadTimeout:         s.config.ReadTimeout,
		ReadHeaderTimeout:   s.config.ReadHeaderTimeout,
		WriteTimeout:        s.config.WriteTimeout,
		IdleTimeout:         s.config.IdleTimeout,
		MaxHeaderBytes:      s.config.MaxHeaderBytes,
		MaxHeaderValueCount: s.config.MaxHeaderValueCount,
		ErrorLog:            log.New(newErrorLogger(s.logger), "", log.LstdFlags|log.Lmicroseconds|log.Llongfile),
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	shutdown := make(chan struct{})
	go func() {
		sig := <-c
		s.logInfo(fmt.Sprintf("teahttp: shutdown: signal: %v", sig))
		ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			s.logError(fmt.Sprintf("teahttp: shutdown: error: %v", err))
		}
		close(shutdown)
	}()
	s.logInfo(fmt.Sprintf("teahttp: listening: %v", server.Addr))
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.logError(fmt.Sprintf("teahttp: error: %v", err))
	}
	<-shutdown
	s.logInfo("teahttp: shutdown")
}

func newErrorLogger(logger Logger) *errorLogger {
	return &errorLogger{logger}
}

type errorLogger struct {
	logger Logger
}

func (w *errorLogger) Write(p []byte) (n int, err error) {
	w.logger.Error(context.Background(), string(p))
	return len(p), nil
}

func (s *Server) logInfo(msg string) {
	s.logger.Info(context.Background(), msg)
	fmt.Fprintln(os.Stdout, msg)
}

func (s *Server) logError(msg string) {
	s.logger.Error(context.Background(), msg)
	fmt.Fprintln(os.Stderr, msg)
}
