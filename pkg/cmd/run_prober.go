package cmd

import (
	"fmt"
	"net/url"
	"os"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"log"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	prober_v1 "kmodules.xyz/prober/api/v1"
	"kmodules.xyz/prober/probe"
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

			return RunProbes(config)
		},
	}
	return cmd
}

func RunProbes(config *rest.Config) error {

	probes := []prober_v1.Handler{
		{
			HTTPGet: &v1.HTTPGetAction{
				Path:        "/success",
				Port:        intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
				Host:        "127.0.0.1",
				Scheme:      "HTTP",
				HTTPHeaders: nil,
			},
		},
		{
			HTTPGet: &v1.HTTPGetAction{
				Path:        "/fail",
				Port:        intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
				Host:        "127.0.0.1",
				Scheme:      "HTTP",
				HTTPHeaders: nil,
			},
		},
		{
			HTTPPost: &prober_v1.HTTPPostAction{
				Path:        "/post-demo",
				Port:        intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
				Host:        "127.0.0.1",
				Scheme:      "HTTP",
				HTTPHeaders: nil,
				Body:        `{"expectedCode":"200","expectedResponse":"success"}`,
			},
		},
		{
			HTTPPost: &prober_v1.HTTPPostAction{
				Path:        "/post-demo",
				Port:        intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
				Host:        "127.0.0.1",
				Scheme:      "HTTP",
				HTTPHeaders: nil,
				Body:        `{"expectedCode":"400","expectedResponse":"failure"}`,
			},
		},
		{
			HTTPPost: &prober_v1.HTTPPostAction{
				Path:        "/post-demo",
				Port:        intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
				Host:        "127.0.0.1",
				Scheme:      "HTTP",
				HTTPHeaders: nil,
				Form: &url.Values{
					"expectedResponse": {"success"},
					"expectedCode":     {"202"},
				},
			},
		},
		{
			HTTPPost: &prober_v1.HTTPPostAction{
				Path:        "/post-demo",
				Port:        intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
				Host:        "127.0.0.1",
				Scheme:      "HTTP",
				HTTPHeaders: nil,
				Form: &url.Values{
					"expectedResponse": {"failure"},
					"expectedCode":     {"404"},
				},
			},
		},
		{
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.IntOrString{Type: intstr.Int, IntVal: 9090},
				Host: "127.0.0.1",
			},
		},
		{
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.IntOrString{Type: intstr.Int, IntVal: 9091},
				Host: "127.0.0.1",
			},
		},
		{
			Exec: &v1.ExecAction{
				Command: []string{"/bin/sh", "-c", `exit $EXIT_CODE_SUCCESS`},
			},
		},
		{
			Exec: &v1.ExecAction{
				Command: []string{"/bin/sh", "-c", `exit $EXIT_CODE_FAIL`},
			},
		},
	}

	kubeClient := kubernetes.NewForConfigOrDie(config)
	pod, err := kubeClient.CoreV1().Pods("default").Get("prober-demo", metav1.GetOptions{})
	if err != nil {
		return err
	}
	status := pod.Status
	container := pod.Spec.Containers[0]

	pb := probe.NewProber(config)

	for i := range probes {
		fmt.Printf("============== Probe No: %d =================\n", i)
		result, ss, err := pb.RunProbe(&probes[i], pod, status, container, time.Second*30)
		if err != nil {
			return err
		}
		fmt.Printf("Result: %v\nReason: %v\n", result, ss)
	}

	return nil
}
