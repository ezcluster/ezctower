package config

import (
	"github.com/go-logr/logr"
)

var (
	Conf Config
	Log  logr.Logger
)

type LogConfig struct {
	Level string `yaml:"level"`
	Mode  string `yaml:"mode"`
}

type Config struct {
	Log     LogConfig `yaml:"log"`
	Workdir string    `yaml:"workdir"`
	Repo    string    `yaml:"repo"`
	Branch  string    `yaml:"branch"`
	Path    string    `yaml:"path"`
}
