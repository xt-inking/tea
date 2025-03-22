package teaconfig

import (
	"time"

	"github.com/BurntSushi/toml"
)

var Config = config{
	Server: server{
		Address:           ":8080",
		ReadTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Second * 0,
		WriteTimeout:      time.Second * 20,
		IdleTimeout:       time.Second * 60,
		MaxHeaderBytes:    1 << 20,
		ShutdownTimeout:   time.Second * 60,
	},
}

type config struct {
	Server server
}

type server struct {
	Address           string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
	ShutdownTimeout   time.Duration
}

func Decode() {
	_, err := toml.DecodeFile(filePath, &Config)
	if err != nil {
		panic(err)
	}
}

const filePath = "manifest/config/config.toml"
