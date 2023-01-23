package config

import (
	"github.com/go-logr/logr"
)

const (
	Marker = ".ezctmarker"
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
	Log          LogConfig `yaml:"log"`
	Auth         Auth      `yaml:"auth"`
	Workdir      string    `yaml:"workdir"`
	RepoUrl      string    `yaml:"repo"` // The git repo URL
	Branch       string    `yaml:"branch"`
	LocalPath    string    `yaml:"path"` // The current path, relative to RepoBasePath
	Committer    Committer `yaml:"committer"`
	RepoName     string    // Computed
	RepoBasePath string    // Computed. RepoBasePath repo location (workdir + repo name)
}
