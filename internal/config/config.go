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

type Committer struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type Auth struct {
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
}

type Config struct {
	Log       LogConfig `yaml:"log"`
	Auth      Auth      `yaml:"auth"`
	Workdir   string    `yaml:"workdir"`
	Repo      string    `yaml:"repo"`
	Branch    string    `yaml:"branch"`
	Path      string    `yaml:"path"`
	Committer Committer `yaml:"committer"`
}
