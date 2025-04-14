package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(releaseCmd)
	rootCmd.AddCommand(recommendCmd)
}

var rootCmd = &cobra.Command{
	Use:   "kube-score",
	Short: "kube-score is the kubernetes release scorer",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
