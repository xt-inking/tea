package teaserver

import (
	"github.com/tea-frame-go/tea/tealog"
)

var ServerOptions = serverOptions{}

type serverOptions struct{}

func (serverOptions) Logger(logger *tealog.Logger) serverOption {
	return func(s *Server) {
		s.logger.Close()
		s.logger = logger
	}
}

type serverOption func(s *Server)
