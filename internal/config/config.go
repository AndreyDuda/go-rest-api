package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"rest_api/pkg/logging"
	"sync"
)

const pathToConfigFile = "config.yml"

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-default:"false"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}

		if err := cleanenv.ReadConfig(pathToConfigFile, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})

	return instance
}
