package cmd

import (
	"ezcluster/tower/internal/config"
	"ezcluster/tower/internal/gitrepo"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

func init() {
	rootCmd.AddCommand(markAsCleanCmd)
}

var markAsCleanCmd = &cobra.Command{
	Use:   "markasclean",
	Short: "Mark the repo as 'clean' (Git content and effective deployment are in sync)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Mark as Clean")
		if err := markAsClean(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error on mark as clean: %v\n", err)
			os.Exit(2)
		}
	},
}

func markAsClean() error {
	gr, err := gitrepo.New(os.Stdout)
	if err != nil {
		return err
	}
	absPath := gr.AbsPath(config.Conf.Path)
	content := time.Now().String()
	err = os.WriteFile(filepath.Join(absPath, config.Marker), []byte(content+"\n"), 0644)
	if err != nil {
		return err
	}
	err = gr.Add(filepath.Join(config.Conf.Path, config.Marker))
	if err != nil {
		return fmt.Errorf("error pn Add(%s): %w", filepath.Join(config.Conf.Path, config.Marker), err)
	}
	err = gr.Commit(fmt.Sprintf("Ecluster tower marker '%s'", content))
	if err != nil {
		return fmt.Errorf("error on Commit(): %w", err)
	}
	err = gr.Push()
	if err != nil {
		return fmt.Errorf("error on Push(): %w", err)
	}
	return nil
}
