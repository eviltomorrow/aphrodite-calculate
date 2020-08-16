package config

import (
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
)

// Config config
type Config struct {
	Log Log `json:"log" toml:"log"`
}

// Log 日志配置项
type Log struct {
	DisableTimestamp bool   `json:"disable-timestamp" toml:"disable-timestamp"`
	Level            string `json:"level" toml:"level"`
	Format           string `json:"format" toml:"format"`
	FileName         string `json:"filename" toml:"filename"`
	MaxSize          int    `json:"maxsize" toml:"maxsize"`
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
	Log: Log{
		DisableTimestamp: false,
		Level:            "info",
		Format:           "text",
		FileName:         "/tmp/aphrodite-calculate/data.log",
		MaxSize:          200,
	},
}
