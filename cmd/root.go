package cmd

import (
	"ezcluster/tower/internal/config"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "tower",
	Short: "A tool aimed to automate ezcluster deployment",
	//Run: func(cmd *cobra.Command, args []string) {
	//	config.Log.V(0).Info("Log V0")
	//	config.Log.V(1).Info("Log V1")
	//	config.Log.V(2).Info("Log V2")
	//	config.Log.Error(errors.New("just to test error message"), "Test ERROR")
	//},
}

func init() {
	config.InitConfig(rootCmd)
	rootCmd.PersistentPreRun = func(command *cobra.Command, args []string) {
		if command != versionCmd {
			if err := config.Setup(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Unable to initialize configuration: %v\n", err)
				os.Exit(2)
			}
		}
	}
}

var debug = true

func Execute() {
	defer func() {
		if !debug {
			if r := recover(); r != nil {
				fmt.Printf("ERROR:%v\n", r)
				os.Exit(1)
			}
		}
	}()
	if err := rootCmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(2)
	}
}
