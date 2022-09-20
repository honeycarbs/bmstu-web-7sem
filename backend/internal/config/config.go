package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"neatly/pkg/logging"
	"sync"
)

type Listen struct {
	Type   string `yaml:"type" env-default:"port"`
	BindIP string `yaml:"bind_ip" env-default:"localhost"`
	Port   string `yaml:"port" env-default:"8080"`
}

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  Listen
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Reading application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("meta/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})

	return instance
}
