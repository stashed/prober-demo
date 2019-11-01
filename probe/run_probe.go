package probe

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
)

func RunProbes(kubeClient kubernetes.Interface) error {
	refManager := kubecontainer.NewRefManager()
	recorder := &record.FakeRecorder{}

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
					Command: []string{"echo", "hello world"},
				},
			},
		},
	}

	pod, err := kubeClient.CoreV1().Pods("default").Get("prober-demo", metav1.GetOptions{})
	if err != nil {
		return err
	}
	status := pod.Status
	container := pod.Spec.Containers[0]

	containerStatus, ok := podutil.GetContainerStatus(status.ContainerStatuses, container.Name)
	if !ok || len(containerStatus.ContainerID) == 0 {
		return fmt.Errorf("failed to extract containerStatus")
	}

	containerID := kubecontainer.ParseContainerID(containerStatus.ContainerID)

	pb := newProber(nil, refManager, recorder)

	for i := range probes {
		result, ss, err := pb.runProbe(readiness, &probes[i], pod, status, container, containerID)
		if err != nil {
			return err
		}
		fmt.Printf("============== Probe No: %d =================\n", i)
		fmt.Printf("Result: %v\nReason: %v\n", result, ss)
	}

	return nil
}
