// Package config TODO
package config

import (
	"time"

	"github.com/qianlnk/config"
)

// Config TODO
type Config struct {
	DownloadFrequency time.Duration `yaml:"DownloadFrequency"`
	GitPushFrequency  time.Duration `yaml:"GitPushFrequency"`
	MysqlDSN          string        `yaml:"MysqlDSN"`
	MarkdownPath      string        `yaml:"MarkdownPath"`
	SourceImagePath   string        `yaml:"SourceImagePath"`
	Owner             string        `yaml:"Owner"`
	MustCommit        int           `yaml:"MustCommit"`
}

// Load TODO
func Load(path string) (*Config, error) {
	cfg := new(Config)
	err := config.Parse(cfg, path)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
