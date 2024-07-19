package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env" env:"ENV" env-default:"local"`
	Grpc     GrpcConfig     `yaml:"grpc"          env-required:"true"`
	Postgres PostgresConfig `yaml:"postgres"      env-required:"true"`
}

type GrpcConfig struct {
	Host    string        `yaml:"host"    env:"GRPC_HOST"    env-default:"localhost"`
	Port    int           `yaml:"port"    env:"GRPC_PORT"    env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT" env-required:"true"`
}

type PostgresConfig struct {
	Host   string `yaml:"host" env:"PG_HOST" env-default:"localhost"`
	Port   int    `yaml:"port" env:"PG_PORT" env-required:"true"`
	User   string `yaml:"user" env:"PG_USER" env-required:"true"`
	Pass   string `yaml:"pass" env:"PG_PASS" env-required:"true"`
	DbName string `yaml:"db"   env:"PG_DB"   env-required:"true"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
