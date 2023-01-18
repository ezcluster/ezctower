package cmd

import (
	"ezcluster/tower/internal/config"
	"ezcluster/tower/internal/gitrepo"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(refreshCmd)
}

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Pull latest version from git server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Refresh")
		if err := refresh(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error on refresh: %v\n", err)
			os.Exit(2)
		}
	},
}

func refresh() error {
	gr, err := gitrepo.New(os.Stdout)
	if err != nil {
		return err
	}
	return refresh2(gr)
}

func refresh2(gr *gitrepo.GitRepo) error {
	b, err := gr.Pull()
	if err != nil {
		return err
	}
	if b {
		config.Log.Info("repo was updated")
	} else {
		config.Log.Info("repo was up-to-date")
	}
	return nil
}
