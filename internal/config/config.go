package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	StoragePath     string        `yaml:"storage_path" env-required:"true"`
	SecretKey       string        `yaml:"secret_key" env-required:"true"`
	ExpirationTime  time.Duration `yaml:"expiration_time" env-default:"1m"`
	AdminCredential string        `yaml:"admin_credential" env-default:"admin"`
	Salt            string        `yaml:"salt" env-default:"yeeeeaaaaaahhSAAALLT"`
	Databases       []Database    `yaml:"databases"`
	HTTPServer      `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Database struct {
	Name             string `yaml:"name"`
	ConnectionString string `yaml:"connection_string"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("CONFIG_PATH environment variable is not set")
	}
	if _, err := os.Stat(configPath); err != nil {
		panic(fmt.Errorf("error opening config file: %s", err))
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(fmt.Errorf("error reading config file: %s", err))
	}

	return &cfg
}

func (c *Config) GetConnectionString(dbName string) (string, error) {
	for _, db := range c.Databases {
		if db.Name == dbName {
			return db.ConnectionString, nil
		}
	}
	return "", errors.New("there is no connections strings for this database")
}
