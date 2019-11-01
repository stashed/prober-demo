package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tamalsaha/prober-demo/probe"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"log"
)

func NewCmdRunProbe() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run-probe",
		Short: "run probe",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Running... probe")
			kubeconfigPath := os.Getenv("KUBECONFIG")
			config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
			if err != nil {
				log.Fatalf("Could not get Kubernetes config: %s", err)
			}
			kubeClient := kubernetes.NewForConfigOrDie(config)
			return probe.RunProbes(kubeClient)
		},
	}
	return cmd
}
