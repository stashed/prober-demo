package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tamalsaha/prober-demo/probe"
)

func NewCmdRunProbe() *cobra.Command  {
	cmd := &cobra.Command{
		Use:   "run-probe",
		Short: "run probe",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Running... probe")
			return probe.RunProbes()
		} ,
	}
	return cmd
}