package session

import (
	"github.com/ilyakaznacheev/cleanenv"
	"neatly/pkg/logging"
	"os"
	"sync"
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

type Swagger struct {
	Host string `yaml:"host"`
}

type Config struct {
	IsDebug *bool   `yaml:"is_debug"`
	DB      DB      `yaml:"db"`
	Listen  Listen  `yaml:"listen"`
	JWT     JWT     `yaml:"jwt"`
	Swagger Swagger `yaml:"swagger"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Infof("Reading application config from %v", os.Getenv("ENV_FILE"))
		instance = &Config{}
		if err := cleanenv.ReadConfig(os.Getenv("ENV_FILE"), instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
