package teaconfig

import (
	"time"

	"github.com/BurntSushi/toml"
)

func NewServerConfig() ServerConfig {
	config := ServerConfig{
		Address:           ":8080",
		ReadTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Second * 0,
		WriteTimeout:      time.Second * 20,
		IdleTimeout:       time.Second * 60,
		MaxHeaderBytes:    1 << 20,
		ShutdownTimeout:   time.Second * 60,
	}
	return config
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

func NewCorsConfig() CorsConfig {
	config := CorsConfig{
		AllowOrigins: []string{},
		AllowHeaders: "Content-Type",
	}
	return config
}

type CorsConfig struct {
	AllowOrigins []string
	AllowHeaders string
}

func NewSqlConfig() SqlConfig {
	config := SqlConfig{
		DriverName:      "mysql",
		DataSourceName:  "root:password@tcp(localhost:3306)/tea-skeleton?charset=utf8mb4&interpolateParams=true&loc=Local&parseTime=true",
		MaxOpenConns:    100,
		MaxIdleConns:    100,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}
	return config
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
