package config

import (
	"flag"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RestConf     RestAPIConfig     `yaml:"restapi"`
	GrpcConf     GRPCConfig        `yaml:"grpc"`
	ConnConf     ConnectionsConfig `yaml:"connections"`
	PostgresConf PostgresConfig    `yaml:"postgres"`
	LogConf      LoggerConfig      `yaml:"logger"`
	RedisConf    RedisConfig       `yaml:"redis"`
}

type RestAPIConfig struct {
	Port    uint32        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type GRPCConfig struct {
	Port    uint32        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type ConnectionsConfig struct {
	ProjServiceConf ProjectServiceConfig `yaml:"projectservice"`
	TaskServiceConf TaskServiceConfig    `yaml:"taskserivce"`
}

type ProjectServiceConfig struct {
	Port    uint32        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type TaskServiceConfig struct {
	Port    uint32        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     uint32 `yaml:"port"`
	User     string `yaml:"user"`
	DbName   string `yaml:"db_name"`
	Password string `yaml:"password"`
	Sslmode  string `yaml:"sslmode"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type RedisConfig struct {
}

func MustLoad() Config {
	confPath := fetchConfigPath()

	if confPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(confPath); err != nil {
		panic("cannot open config path")
	}

	data, err := os.ReadFile(confPath)
	if err != nil {
		panic("cannot read config path")
	}

	var config Config

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		panic("cannot parse config path")
	}

	return config
}

func fetchConfigPath() string {
	var confPath string

	flag.StringVar(&confPath, "config", "", "path to config")
	flag.Parse()

	if confPath == "" {
		confPath = os.Getenv("CONFIG")
	}

	return confPath
}
