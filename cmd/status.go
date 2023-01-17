package cmd

import (
	"ezcluster/tower/internal/config"
	"ezcluster/tower/internal/gitrepo"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get status (Dirty or clean)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Status")
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
	hash1, err := gr.GetLastHashLog(config.Conf.Path)
	if err != nil {
		return false, err
	}
	hash2, err := gr.GetLastHashLog(filepath.Join(config.Conf.Path, config.Marker))
	if err != nil {
		return false, err
	}
	config.Log.V(1).Info("maker log", "hash base", hash1, "hash marker", hash2)

	return hash1 == hash2, nil
}
