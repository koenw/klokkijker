package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	GitCommit string

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  `Show version information`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("klokkijker", GitCommit)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
