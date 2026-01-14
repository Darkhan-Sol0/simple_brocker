package config

import (
	"log"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	config struct {
		Address  string                 `yaml:"address"`
		Services map[string]serviceConf `yaml:"services"`
	}

	serviceConf struct {
		Address   string        `yaml:"address"`
		CoolDown  time.Duration `yaml:"cooldown"`
		BatchSize int           `yaml:"batch_size"`
	}

	ServiceConf interface {
		GetServiceAddress() string
		GetServiceBatchSize() int
		GetCoolDown() time.Duration
	}

	Config interface {
		GetAddress() string
		GetServices() map[string]serviceConf
		GetService(key string) ServiceConf
	}
)

var (
	conf Config
	once sync.Once
)

func GetConfig() Config {
	once.Do(func() {
		conf = &config{}
		path := "./config/config.yaml"
		if err := cleanenv.ReadConfig(path, conf); err != nil {
			log.Fatalf("error read config file %s: %v", path, err)
		}
	})
	return conf
}

func (c *config) GetAddress() string {
	return c.Address
}

func (c *serviceConf) GetCoolDown() time.Duration {
	return c.CoolDown
}

func (c *config) GetServices() map[string]serviceConf {
	return c.Services
}

func (c *config) GetService(key string) ServiceConf {
	s := c.Services[key]
	return &s
}

func (c *serviceConf) GetServiceAddress() string {
	return c.Address
}

func (c *serviceConf) GetServiceBatchSize() int {
	return c.BatchSize
}
