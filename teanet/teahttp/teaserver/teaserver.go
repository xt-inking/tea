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
}

func New(options ...serverOption) *Server {
	s := &Server{
		router: newRouter(),
		logger: tealog.New(
			tealog.NewRecordHandlerText(),
			tealog.NewWriterCloserFile(tealog.FileDir, "http"),
		),
		loggers: make([]*tealog.Logger, 0),
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
		Addr:                         teaconfig.Server.Address,
		Handler:                      http.HandlerFunc(s.handler),
		DisableGeneralOptionsHandler: true,
		ReadTimeout:                  teaconfig.Server.ReadTimeout,
		ReadHeaderTimeout:            teaconfig.Server.ReadHeaderTimeout,
		WriteTimeout:                 teaconfig.Server.WriteTimeout,
		IdleTimeout:                  teaconfig.Server.IdleTimeout,
		MaxHeaderBytes:               teaconfig.Server.MaxHeaderBytes,
		ErrorLog:                     log.New(newErrorLogger(s.logger), "", log.LstdFlags|log.Lmicroseconds|log.Llongfile),
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	closed := make(chan struct{})
	go func() {
		sig := <-c
		s.logger.WithStdout().Infof(
			context.Background(),
			"http: shutting down by signal `%s`",
			sig,
		)
		ctx, cancel := context.WithTimeout(context.Background(), teaconfig.Server.ShutdownTimeout)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			s.logger.WithStdout().Errorf(
				context.Background(),
				"http: shutdown error: %v",
				err,
			)
		}
		close(closed)
	}()
	s.logger.WithStdout().Infof(
		context.Background(),
		"http: listening on `%s`",
		server.Addr,
	)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.logger.WithStdout().Errorf(
			context.Background(),
			"http: error: %v",
			err,
		)
	}
	<-closed
	s.logger.WithStdout().Infof(
		context.Background(),
		"http: shutdown",
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
