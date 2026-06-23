package teaserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tea-frame-go/tea/teaconfig"
	"github.com/tea-frame-go/tea/tealog"
)

type Server struct {
	router  *router
	logger  *tealog.Logger
	loggers []*tealog.Logger
	config  *teaconfig.ServerConfig
}

func New(config *teaconfig.ServerConfig, options ...serverOption) *Server {
	s := &Server{
		router: newRouter(),
		logger: tealog.New(
			tealog.NewRecordText,
			tealog.NewWriterCloserFile(tealog.FileDir, "http"),
		),
		loggers: make([]*tealog.Logger, 0),
		config:  config,
	}
	for _, o := range options {
		o(s)
	}
	return s
}

func (s *Server) Router() *router {
	return s.router
}

func (s *Server) Loggers(loggers ...*tealog.Logger) {
	s.loggers = append(s.loggers, loggers...)
}

func (s *Server) Serve() {
	defer func() {
		s.logger.Close()
		for i := range s.loggers {
			s.loggers[i].Close()
		}
	}()
	server := &http.Server{
		Addr:                         s.config.Address,
		Handler:                      http.HandlerFunc(s.handler),
		DisableGeneralOptionsHandler: true,
		ReadTimeout:                  s.config.ReadTimeout,
		ReadHeaderTimeout:            s.config.ReadHeaderTimeout,
		WriteTimeout:                 s.config.WriteTimeout,
		IdleTimeout:                  s.config.IdleTimeout,
		MaxHeaderBytes:               s.config.MaxHeaderBytes,
		ErrorLog:                     log.New(newErrorLogger(s.logger), "", log.LstdFlags|log.Lmicroseconds|log.Llongfile),
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	closed := make(chan struct{})
	go func() {
		sig := <-c
		s.logger.WithStdout().Infof(
			context.Background(),
			"teahttp: shutting down by signal `%s`",
			sig,
		)
		ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			s.logger.WithStdout().Errorf(
				context.Background(),
				"teahttp: shutdown error: %v",
				err,
			)
		}
		close(closed)
	}()
	s.logger.WithStdout().Infof(
		context.Background(),
		"teahttp: listening on `%s`",
		server.Addr,
	)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.logger.WithStdout().Errorf(
			context.Background(),
			"teahttp: error: %v",
			err,
		)
	}
	<-closed
	s.logger.WithStdout().Info(
		context.Background(),
		"teahttp: shutdown",
	)
}

type errorLogger struct {
	logger *tealog.Logger
}

func newErrorLogger(logger *tealog.Logger) *errorLogger {
	return &errorLogger{
		logger: logger,
	}
}

func (e *errorLogger) Write(p []byte) (n int, err error) {
	e.logger.Error(context.Background(), string(p))
	return 0, nil
}
