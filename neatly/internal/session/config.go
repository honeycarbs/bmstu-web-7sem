package session

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

type JWTConfig struct {
	Secret string `yaml:"secret" env-required:"true"`
}

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  Listen
	JWT     JWTConfig
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Reading application session")
		instance = &Config{}
		if err := cleanenv.ReadConfig("meta/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})

	return instance
}
