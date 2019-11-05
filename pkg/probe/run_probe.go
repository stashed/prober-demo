package probe

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"k8s.io/apimachinery/pkg/util/intstr"
	v1 "stash.appscode.dev/prober-demo/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func RunProbes(config *rest.Config) error {

	probes := []v1.Probe{
		{
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Path:        "/success",
					Port:        intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
					Host:        "127.0.0.1",
					Scheme:      "HTTP",
					HTTPHeaders: nil,
				},
			},
		},
		{
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Path:        "/fail",
					Port:        intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
					Host:        "127.0.0.1",
					Scheme:      "HTTP",
					HTTPHeaders: nil,
				},
			},
		},
		{
			Handler: v1.Handler{
				TCPSocket: &v1.TCPSocketAction{
					Port: intstr.IntOrString{Type: intstr.Int, IntVal: 9090},
					Host: "127.0.0.1",
				},
			},
		},
		{
			Handler: v1.Handler{
				TCPSocket: &v1.TCPSocketAction{
					Port: intstr.IntOrString{Type: intstr.Int, IntVal: 9091},
					Host: "127.0.0.1",
				},
			},
		},
		{
			Handler: v1.Handler{
				Exec: &v1.ExecAction{
					Command: []string{"/bin/sh", "-c",`exit $EXIT_CODE_SUCCESS`},
				},
			},
		},
		{
			Handler: v1.Handler{
				Exec: &v1.ExecAction{
					Command: []string{"/bin/sh", "-c",`exit $EXIT_CODE_FAIL`},
				},
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

	pb := newProber(config)

	for i := range probes {
		fmt.Printf("============== Probe No: %d =================\n", i)
		result, ss, err := pb.runProbe(&probes[i], pod, status, container)
		if err != nil {
			return err
		}
		fmt.Printf("Result: %v\nReason: %v\n", result, ss)
	}

	return nil
}
