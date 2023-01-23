package cmd

import (
	"ezcluster/tower/internal/config"
	"fmt"
	"github.com/spf13/cobra"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(locationCmd)
	locationCmd.AddCommand(baseCmd)
	locationCmd.AddCommand(targetCmd)
	locationCmd.AddCommand(reponameCmd)
}

var locationCmd = &cobra.Command{
	Use:   "location",
	Short: "Display repo location",
}

var baseCmd = &cobra.Command{
	Use:   "base",
	Short: "Display base repo location",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Conf.RepoBasePath)
	},
}

var targetCmd = &cobra.Command{
	Use:   "target",
	Short: "Display target (cluster) folder location",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(filepath.Join(config.Conf.RepoBasePath, config.Conf.LocalPath))
	},
}

var reponameCmd = &cobra.Command{
	Use:   "reponame",
	Short: "Display repo name",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Conf.RepoName)
	},
}
