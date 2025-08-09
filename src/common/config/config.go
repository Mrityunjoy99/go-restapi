package config

import "github.com/creasty/defaults"

type Config struct {
	App AppConfig
	DB  DatabaseConfig
	Jwt JwtConfig
}

type AppConfig struct {
	Name string `default:"sample-go"`
	Port int    `default:"8080"`
	Id   string `default:"sample-go"`
}

type DatabaseConfig struct {
	Host               string `default:"localhost"`
	Port               int    `default:"5432"`
	Name               string `default:""`
	User               string `default:""`
	Password           string `default:""`
	MaxIdleConnections int    `default:"2"`
	MaxOpenConnections int    `default:"5"`
	ConnMaxIdleTimeSec int    `default:"10"`
	ConnMaxLifeTimeSec int    `default:"600"`
}

type JwtConfig struct {
	Secret        string `default:""`
	ExpireTimeSec int    `default:"3600"`
}

var c *Config

func NewConfig() (*Config, error) {
	if c != nil {
		return c, nil
	}

	defaultConfig := &Config{}
	if err := defaults.Set(defaultConfig); err != nil {
		return nil, err
	}

	c, err := LoadConfig(defaultConfig)
	if err != nil {
		return nil, err
	}

	return c, nil
}
