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

var Sql = SqlConfig{
	DriverName:      "mysql",
	DataSourceName:  "username:password@protocol(address)/dbname?param=value",
	MaxOpenConns:    100,
	MaxIdleConns:    100,
	ConnMaxLifetime: 5 * time.Minute,
	ConnMaxIdleTime: 5 * time.Minute,
}

type SqlConfig struct {
	DriverName      string
	DataSourceName  string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func Decode[Config any](config *Config) {
	_, err := toml.DecodeFile(filePath, config)
	if err != nil {
		panic(err)
	}
}

const filePath = "manifest/config/config.toml"
