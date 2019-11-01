package probe

import (
	"fmt"

	"github.com/the-redback/go-oneliners"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
)

func RunProbes() error {
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
	}

	pod := &v1.Pod{}
	status := v1.PodStatus{}
	container := v1.Container{}
	containerID := kubecontainer.ContainerID{}

	pb := newProber(nil, refManager, recorder)

	for i := range probes {
		result, ss, err := pb.runProbe(readiness, &probes[i], pod, status, container, containerID)
		if err != nil {
			return err
		}
		fmt.Println("ss: ", ss)
		oneliners.PrettyJson(result, "result")
	}

	return nil
}
