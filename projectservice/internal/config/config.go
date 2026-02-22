package config

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-yaml"
)

var (
	dockerType = "docker"
	localType  = "local"
)

type Config struct {
	Type            string            `yaml:"type"`
	RestConf        RestAPIConfig     `yaml:"restapi"`
	GrpcConf        GRPCConfig        `yaml:"grpc"`
	ConnectionsConf ConnectionsConfig `yaml:"connections"`
	PostgresConf    PostgresConfig    `yaml:"postgres"`
	LoggerConf      LoggerConfig      `yaml:"logger"`
}

type RestAPIConfig struct {
	Port            uint32        `yaml:"port"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	RequestTimeout  time.Duration `yaml:"request_timeout"`
	Mode            string        `yaml:"mode"`
}

type GRPCConfig struct {
	Host    string        `yaml:"host"`
	Port    uint32        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type ConnectionsConfig struct {
	UserServConnConf UserServiceConnectionConfig `yaml:"userservice"`
}

type UserServiceConnectionConfig struct {
	Host            string        `yaml:"host"`
	Port            uint32        `yaml:"port"`
	ResponseTimeout time.Duration `yaml:"response_timeout"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     uint32 `yaml:"port"`
	User     string `yaml:"user"`
	DbName   string `yaml:"db_name"`
	Password string
	Sslmode  string `yaml:"sslmode"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

func MustLoad() *Config {
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

	loadSecrets(&config)

	return &config
}

func loadSecrets(cfg *Config) {
	if cfg.Type == localType {
		cfg.PostgresConf.Password = os.Getenv("DB_PASS")
		if cfg.PostgresConf.Password == "" {
			panic("PostgresConfig password field empty, (DB_PASS)")
		}
	} else if cfg.Type == dockerType {
		cfg.PostgresConf.Host = os.Getenv("DB_HOST")
		if cfg.PostgresConf.Host == "" {
			panic("PostgresConfig host field empty, (DB_HOST)")
		}
		port, err := strconv.Atoi(os.Getenv("DB_PORT"))
		if err != nil || port == 0 {
			panic("PostgresConf port error (DB_PORT): " + err.Error())
		}
		cfg.PostgresConf.Port = uint32(port)
		cfg.PostgresConf.User = os.Getenv("DB_USER")
		if cfg.PostgresConf.User == "" {
			panic("PostgresConf user field empty (DB_USER)")
		}
		cfg.PostgresConf.DbName = os.Getenv("DB_NAME")
		if cfg.PostgresConf.DbName == "" {
			panic("PostgresConf db name file empty, (DB_NAME)")
		}
		cfg.PostgresConf.Password = os.Getenv("DB_PASS")
		if cfg.PostgresConf.Password == "" {
			panic("PostgresConf password field empty, (DB_PASS)")
		}
		cfg.PostgresConf.Sslmode = os.Getenv("DB_MODE")
		if cfg.PostgresConf.Sslmode == "" {
			panic("PostgresConf sslmode field empty, (DB_MODE)")
		}
	}
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
