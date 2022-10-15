package session

import (
	"github.com/ilyakaznacheev/cleanenv"
	"neatly/pkg/logging"
	"sync"
)

const (
	confPath = "etc/config/config.yml"
)

type DB struct {
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	DBName         string `yaml:"dbname"`
	SSLMode        string `yaml:"ssl_mode"`
	MigrationsPath string `yaml:"migrations_path"`
}

type Listen struct {
	Port   string `yaml:"port"`
	BindIP string `yaml:"bind_ip"`
}

type JWT struct {
	Secret string `yaml:"secret"`
}

type Config struct {
	IsDebug *bool  `yaml:"is_debug"`
	DB      DB     `yaml:"db"`
	Listen  Listen `yaml:"listen"`
	JWT     JWT    `yaml:"jwt"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Reading application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig(confPath, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
