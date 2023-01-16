package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(refreshCmd)
}

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Pull latest version from git server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Refresh")
	},
}
