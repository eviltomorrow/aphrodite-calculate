package config

import (
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
)

// Config config
type Config struct {
	MongoDB MongoDB `json:"mongodb" toml:"mongodb"`
	Log     Log     `json:"log" toml:"log"`
	MySQL   MySQL   `json:"mysql" toml:"mysql"`
	System  System  `json:"system" toml:"system"`
}

// MongoDB mongodb
type MongoDB struct {
	Database string `json:"database" toml:"database"`
	DSN      string `json:"dsn" toml:"dsn"`
	MinOpen  uint64 `json:"min-open" toml:"min-open"`
	MaxOpen  uint64 `json:"max-open" toml:"max-open"`
}

// MySQL mysql
type MySQL struct {
	DSN     string `json:"dsn" toml:"dsn"`
	MinOpen int    `json:"min-open" toml:"min-open"`
	MaxOpen int    `json:"max-open" toml:"max-open"`
}

// Log 日志配置项
type Log struct {
	DisableTimestamp bool   `json:"disable-timestamp" toml:"disable-timestamp"`
	Level            string `json:"level" toml:"level"`
	Format           string `json:"format" toml:"format"`
	FileName         string `json:"filename" toml:"filename"`
	MaxSize          int    `json:"maxsize" toml:"maxsize"`
}

// System system
type System struct {
	PProfListenPort int    `json:"pprof-listen-port" toml:"pprof-listen-port"`
	BeginDate       string `json:"begin-date" toml:"begin-date"`
}

// Load 加载配置文件
func (cg *Config) Load(path string, f func(*Config)) error {
	_, err := toml.DecodeFile(path, cg)
	if err != nil {
		return err
	}
	f(cg)
	return nil
}

func (cg *Config) String() string {
	buf, err := json.Marshal(cg)
	if err != nil {
		return fmt.Sprintf("Marshal config to json failure, nest error: %v", err)
	}
	return string(buf)
}

// DefaultGlobalConfig default config
var DefaultGlobalConfig = &Config{
	MongoDB: MongoDB{
		Database: "aphrodite",
		DSN:      "mongodb://localhost:27017",
		MinOpen:  5,
		MaxOpen:  10,
	},
	MySQL: MySQL{
		DSN:     "root:root@tcp(localhost:3306)/aphrodite?charset=utf8mb4&parseTime=true&loc=Local",
		MinOpen: 5,
		MaxOpen: 10,
	},
	Log: Log{
		DisableTimestamp: false,
		Level:            "info",
		Format:           "text",
		FileName:         "/tmp/aphrodite-calculate/data.log",
		MaxSize:          20,
	},
	System: System{
		PProfListenPort: 6070,
		BeginDate:       "2020-08-31",
	},
}
