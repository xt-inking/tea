package teaconfig

import (
	"time"

	"github.com/BurntSushi/toml"
)

var Server = ServerConfig{
	Address:           ":8080",
	ReadTimeout:       time.Second * 10,
	ReadHeaderTimeout: time.Second * 0,
	WriteTimeout:      time.Second * 20,
	IdleTimeout:       time.Second * 60,
	MaxHeaderBytes:    1 << 20,
	ShutdownTimeout:   time.Second * 60,
}

type ServerConfig struct {
	Address           string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
	ShutdownTimeout   time.Duration
}

func Decode[Config any](config *Config) {
	_, err := toml.DecodeFile(filePath, config)
	if err != nil {
		panic(err)
	}
}

const filePath = "manifest/config/config.toml"
