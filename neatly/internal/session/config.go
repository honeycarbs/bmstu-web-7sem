package session

import (
	"github.com/ilyakaznacheev/cleanenv"
	"neatly/pkg/logging"
	"sync"
)

type DB struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"ssl_mode"`
}

type Listen struct {
	Port   string `yaml:"port" env-default:"8080"`
	BindIP string `yaml:"bind_ip" env-default:"localhost"`
}

type JWT struct {
	Secret string `yaml:"secret"`
}

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	DB      DB    `yaml:"db"`
	Listen  Listen
	JWT     JWT `yaml:"jwt"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Reading application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
