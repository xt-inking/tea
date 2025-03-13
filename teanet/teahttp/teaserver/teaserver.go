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

type server struct {
	router *router
	logger *tealog.Logger
}

func New() *server {
	s := &server{
		router: newRouter(),
		logger: tealog.New(
			tealog.NewRecordHandlerText(),
			tealog.NewWriterCloserFile(tealog.FileDir, "http"),
		),
	}
	return s
}

func (s *server) Router() *router {
	return s.router
}

func (s *server) Serve() {
	defer func() {
		s.logger.Close()
	}()
	server := &http.Server{
		Addr:                         teaconfig.Config.Server.Address,
		Handler:                      http.HandlerFunc(s.handler),
		DisableGeneralOptionsHandler: true,
		ReadTimeout:                  teaconfig.Config.Server.ReadTimeout,
		ReadHeaderTimeout:            teaconfig.Config.Server.ReadHeaderTimeout,
		WriteTimeout:                 teaconfig.Config.Server.WriteTimeout,
		IdleTimeout:                  teaconfig.Config.Server.IdleTimeout,
		MaxHeaderBytes:               teaconfig.Config.Server.MaxHeaderBytes,
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
		ctx, cancel := context.WithTimeout(context.Background(), teaconfig.Config.Server.ShutdownTimeout)
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
