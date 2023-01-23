package cmd

import (
	"ezcluster/tower/internal/config"
	"ezcluster/tower/internal/gitrepo"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var skipRefresh bool

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.PersistentFlags().BoolVar(&skipRefresh, "skipRefresh", false, "")
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Refresh and get status (Dirty or clean)",
	Run: func(cmd *cobra.Command, args []string) {
		clean, err := status()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error on status: %v\n", err)
			os.Exit(2)
		}
		if clean {
			fmt.Printf("Clean\n")
			os.Exit(0)
		} else {
			fmt.Printf("Dirty\n")
			os.Exit(1)
		}
	},
}

func status() (clean bool, err error) {
	gr, err := gitrepo.New(os.Stdout)
	if err != nil {
		return false, err
	}
	if !skipRefresh {
		err = refresh2(gr)
		if err != nil {
			return false, err
		}
	}
	absPath := gr.AbsPath(config.Conf.LocalPath)
	_, err = os.Stat(filepath.Join(absPath, config.Marker))
	if err != nil {
		config.Log.V(1).Info("marker not set", "path", config.Conf.LocalPath, "error", err.Error())
		// .marker does not exists. Let's say it is dirty, as never makasclean-ed
		return false, nil
	}

	hash1, err := gr.GetLastHashLog(config.Conf.LocalPath)
	if err != nil {
		return false, err
	}
	hash2, err := gr.GetLastHashLog(filepath.Join(config.Conf.LocalPath, config.Marker))
	if err != nil {
		return false, err
	}
	config.Log.V(1).Info("marker log", "hashBase", hash1, "hashMarker", hash2)

	return hash1 == hash2, nil
}
