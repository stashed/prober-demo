package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command  {
	rootCmd := &cobra.Command{
		Use:   "prober",
		Short: "prober root command",
	}
	rootCmd.AddCommand(NewCmdRunProbe())
	rootCmd.AddCommand(NewCmdRunClient())
	return rootCmd
}