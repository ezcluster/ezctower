package config

import (
	"ezcluster/tower/pkg/varenvs"
	"fmt"
	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

var (
	configFile string
	logLevel   string
	logMode    string
)

func InitConfig(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVar(&configFile, "configFile", "", "Configuration file")
	rootCmd.PersistentFlags().StringVar(&logLevel, "logLevel", "", "Log level (PANIC|FATAL|ERROR|WARN|INFO|DEBUG|TRACE)")
	rootCmd.PersistentFlags().StringVar(&logMode, "logMode", "", "Log mode: 'dev' or 'json'")
}

func Setup() error {
	// ------------------------ Load config file, if any
	if configFile == "" {
		Conf = Config{}
	} else {
		absConfigFile, err := filepath.Abs(configFile)
		if err != nil {
			return err
		}
		file, err := os.Open(absConfigFile)
		if err != nil {
			return err
		}
		decoder := yaml.NewDecoder(file)
		decoder.KnownFields(true)
		if err = decoder.Decode(&Conf); err != nil {
			return fmt.Errorf("file '%s': %w", configFile, err)
		}
		// Adjust workdir to configFile path if not absolute
		if !filepath.IsAbs(Conf.Workdir) {
			Conf.Workdir = filepath.Join(filepath.Dir(absConfigFile), Conf.Workdir)
		}
	}
	// ---------------------- Combine with environment variable (Which take precedence)
	varenv := varenvs.New()
	varenv.Add("workdir", &Conf.Workdir, "EZCT_WORKDIR", "", true, false)
	varenv.Add("repo", &Conf.Repo, "EZCT_REPO", "", true, false)
	varenv.Add("branch", &Conf.Branch, "EZCT_BRANCH", "", true, false)
	varenv.Add("path", &Conf.Path, "EZCT_PATH", "", true, false)
	varenv.Add("user", &Conf.Auth.Username, "EZCT_USERNAME", "", false, false)
	varenv.Add("token", &Conf.Auth.Token, "EZCT_TOKEN", "", false, false)
	varenv.Add("committerName", &Conf.Committer.Name, "EZCT_COMMITTER_NAME", "ezctower", false, false)
	varenv.Add("committerEmail", &Conf.Committer.Email, "EZCT_COMMITTER_EMAIL", "tower@ezcluster.com", false, false)

	err := varenv.Parse()
	if err != nil {
		return err
	}
	// -----------------------------------Handle logging  stuff
	Log, err = handleLog(&Conf.Log, logLevel, logMode)
	if err != nil {
		return err
	}
	return nil
}

func handleLog(logConfig *LogConfig, logLevel string, logMode string) (logr.Logger, error) {

	// Override configFile value by command line ones
	if logMode != "" {
		logConfig.Mode = logMode
	}
	if logLevel != "" {
		logConfig.Level = logLevel
	}
	// Set default values
	if logConfig.Mode == "" {
		logConfig.Mode = "json"
	}
	if logConfig.Level == "" {
		logConfig.Level = "INFO"
	}
	logConfig.Mode = strings.ToLower(logConfig.Mode)
	logConfig.Level = strings.ToUpper(logConfig.Level)

	if logConfig.Mode != "dev" && logConfig.Mode != "json" {
		return logr.New(nil), fmt.Errorf("invalid logMode value: %s. Must be one of 'dev' or 'json'", logConfig.Mode)
	}
	llevel, ok := logLevelByString[logConfig.Level]
	if !ok {
		return logr.New(nil), fmt.Errorf("%s is an invalid value for Log.Level\n", logConfig.Level)
	}

	logrusLog := logrus.New()
	logrusLog.SetLevel(llevel)
	if logConfig.Mode == "json" {
		logrusLog.SetFormatter(&logrus.JSONFormatter{})
	}
	l := logrusr.New(logrusLog)
	return l, nil

}

var logLevelByString = map[string]logrus.Level{
	"PANIC": logrus.PanicLevel,
	"FATAL": logrus.FatalLevel,
	"ERROR": logrus.ErrorLevel,
	"WARN":  logrus.WarnLevel,
	"INFO":  logrus.InfoLevel,
	"DEBUG": logrus.DebugLevel,
	"TRACE": logrus.TraceLevel,
}
