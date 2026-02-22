package config

import (
	"flag"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	dockerType = "docker"
	localType  = "local"
)

type Config struct {
	Type         string         `yaml:"type"`
	RestConf     RestAPIConfig  `yaml:"restapi"`
	GrpcConf     GRPCConfig     `yaml:"grpc"`
	PostgresConf PostgresConfig `yaml:"postgres"`
	LogConf      LoggerConfig   `yaml:"logger"`
	RedisConf    RedisConfig    `yaml:"redis"`
}

type RestAPIConfig struct {
	Port            uint32        `yaml:"port"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	RequestTimeout  time.Duration `yaml:"request_timeout"`
	Mode            string        `yaml:"mode"`
	CookieTTL       time.Duration `yaml:"cookie-ttl"`
}

type GRPCConfig struct {
	Host    string        `yaml:"host"`
	Port    uint32        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
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

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     uint32 `yaml:"port"`
	Password string
	DB       int           `yaml:"db"`
	TTL      time.Duration `yaml:"ttl"`
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

	loadSecrets(&config)

	return config
}

func loadSecrets(cfg *Config) {
	if cfg.Type == localType {
		cfg.PostgresConf.Password = os.Getenv("DB_PASS")
		if cfg.PostgresConf.Password == "" {
			panic("PostgresConfig password field empty")
		}
		cfg.RedisConf.Password = os.Getenv("REDIS_PASS")
		if cfg.RedisConf.Password == "" {
			panic("RedisConf password field empty")
		}
	} else if cfg.Type == dockerType {
		mustLoadPostgresConfig(cfg)
		mustLoadRedisConfig(cfg)
	}
}

func mustLoadPostgresConfig(cfg *Config) {
	cfg.PostgresConf.Host = os.Getenv("DB_HOST")
	if cfg.PostgresConf.Host == "" {
		panic("PostgresConfig host field empty")
	}
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	cfg.PostgresConf.Port = uint32(port)
	cfg.PostgresConf.User = os.Getenv("DB_USER")
	if cfg.PostgresConf.User == "" {
		panic("PostgresConfig user field empty")
	}
	cfg.PostgresConf.DbName = os.Getenv("DB_NAME")
	if cfg.PostgresConf.DbName == "" {
		panic("PostgresConfig db name field empty")
	}
	cfg.PostgresConf.Password = os.Getenv("DB_PASS")
	if cfg.PostgresConf.Password == "" {
		panic("PostgresConfig password field empty")
	}
	cfg.PostgresConf.Sslmode = os.Getenv("DB_MODE")
	if cfg.PostgresConf.Sslmode == "" {
		panic("PostgresConfig sslmode field empty")
	}
}

func mustLoadRedisConfig(cfg *Config) {
	cfg.RedisConf.Host = os.Getenv("REDIS_HOST")
	if cfg.RedisConf.Host == "" {
		panic("RedisConf host field empty")
	}
	port, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	cfg.RedisConf.Port = uint32(port)
	cfg.RedisConf.Password = os.Getenv("REDIS_PASS")
	if cfg.RedisConf.Password == "" {
		panic("RedisConf password field empty")
	}
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	cfg.RedisConf.DB = db
	h, _ := strconv.Atoi(os.Getenv("REDIS_TTL_SEC"))
	cfg.RedisConf.TTL = time.Second * time.Duration(h)
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
