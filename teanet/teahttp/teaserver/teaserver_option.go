package teaserver

type option interface {
	apply(s *Server)
}

type optionFunc func(s *Server)

func (f optionFunc) apply(s *Server) {
	f(s)
}

var Opt = opt{}

type opt struct{}

func (opt) Logger(logger Logger) optionFunc {
	return func(s *Server) {
		s.logger.Close()
		s.logger = logger
	}
}

func (opt) Config(f func(config *Config)) optionFunc {
	return func(s *Server) {
		f(s.config)
	}
}
