package config

import (
	"log"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	config struct {
		Address string                `yaml:"address"`
		TLS     TLS                   `yaml:"tls"`
		Groups  map[string]*groupConf `yaml:"group"`
		MaxChan int                   `yaml:"max_chan"`
	}

	TLS struct {
		Enabled  bool   `yaml:"enabled"`
		CertPath string `yaml:"cert_path"`
		KeyPath  string `yaml:"key_path"`
	}

	groupConf struct {
		Address   []string      `yaml:"address"`
		CoolDown  time.Duration `yaml:"cooldown"`
		BatchSize int           `yaml:"batch_size"`
		Retry     int           `yaml:"retry"`
	}

	GroupConf interface {
		GetServiceAddress() []string
		GetServiceBatchSize() int
		GetCoolDown() time.Duration
		GetRetry() int
	}

	Config interface {
		GetTLS() TLS
		GetAddress() string
		GetGroups() map[string]*groupConf
		GetGroup(key string) GroupConf
		CheckGroup(group string) bool
		GetMaxChan() int
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

func (c *config) GetTLS() TLS {
	return c.TLS
}

func (c *config) GetAddress() string {
	return c.Address
}

func (c *config) GetGroups() map[string]*groupConf {
	return c.Groups
}

func (c *config) GetMaxChan() int {
	return c.MaxChan
}

func (c *config) GetGroup(key string) GroupConf {
	return c.Groups[key]

}
func (g *groupConf) GetCoolDown() time.Duration {
	return g.CoolDown
}

func (g *groupConf) GetServiceAddress() []string {
	return g.Address
}

func (g *groupConf) GetServiceBatchSize() int {
	return g.BatchSize
}

func (g *groupConf) GetRetry() int {
	return g.Retry
}

func (c *config) CheckGroup(group string) bool {
	return c.Groups[group] != nil
}
