package config

import (
	"fmt"
	"sync"

	"authkey/pkg/lib/db"
	"authkey/pkg/lib/log"
)

var (
	// GlobalConfig global config
	GlobalConfig *Config
	lock         sync.RWMutex = sync.RWMutex{}
)

// Common Config's Common
type Common struct {
	Env  string
	Name string
	Host string
	Port int
	Desc string
}

func (c *Common) IsDev() bool {
	return c.Env == "dev"
}

// Service Config's Common
type Service struct {
	Host string
	Port int
}

type Trace struct {
	Enable bool
	Server string
}

// Config
type Config struct {
	*Common
	*Service
	*Trace
	Logger *log.ZapConfig
	DB     *db.DBConfig
}

// SetConfig update config
func SetConfig(c *Config) {
	lock.Lock()
	defer lock.Unlock()

	GlobalConfig = c
}

func (c *Config) String() string {
	return fmt.Sprintf("Commont=%+v, Logger=%+v", *c.Common, *c.Logger)
}
