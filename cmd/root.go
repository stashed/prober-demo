package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command  {
	rootCmd := &cobra.Command{
		Use:   "run",
		Short: "run client or prober",
	}
	rootCmd.AddCommand(NewCmdRunProbe())
	rootCmd.AddCommand(NewCmdRunClient())
	return rootCmd
}