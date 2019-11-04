package probe

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	//"k8s.io/client-go/tools/record"
	//kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func RunProbes(kubeClient kubernetes.Interface) error {
	//refManager := kubecontainer.NewRefManager()
	//recorder := &record.FakeRecorder{}

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
		//{
		//	Handler: v1.Handler{
		//		Exec: &v1.ExecAction{
		//			Command: []string{"echo", "hello world"},
		//		},
		//	},
		//},
	}

	fmt.Println(len(probes))
	pod, err := kubeClient.CoreV1().Pods("default").Get("prober-demo", metav1.GetOptions{})
	if err != nil {
		return err
	}
	status := pod.Status
	container := pod.Spec.Containers[0]

	//containerStatus, ok := podutil.GetContainerStatus(status.ContainerStatuses, container.Name)
	//if !ok || len(containerStatus.ContainerID) == 0 {
	//	return fmt.Errorf("failed to extract containerStatus")
	//}

	//containerID := kubecontainer.ParseContainerID(containerStatus.ContainerID)

	pb := newProber()

	for i := range probes {
		result, ss, err := pb.runProbe(&probes[i], pod, status, container)
		if err != nil {
			return err
		}
		fmt.Printf("============== Probe No: %d =================\n", i)
		fmt.Printf("Result: %v\nReason: %v\n", result, ss)
	}

	return nil
}

//// RunInContainer synchronously executes the command in the container, and returns the output.
//func (pb *prober) RunInContainer(id kubecontainer.ContainerID, cmd []string, timeout time.Duration) ([]byte, error) {
//	stdout, stderr, err := pb.ExecSync(id.ID, cmd, timeout)
//	// NOTE(tallclair): This does not correctly interleave stdout & stderr, but should be sufficient
//	// for logging purposes. A combined output option will need to be added to the ExecSyncRequest
//	// if more precise output ordering is ever required.
//	return append(stdout, stderr...), err
//}
//
//func (pb *prober) ExecOnPod(pod *core.Pod, command ...string) (string, error) {
//	var (
//		execOut bytes.Buffer
//		execErr bytes.Buffer
//	)
//
//	req := pb.KubeClient.CoreV1().RESTClient().Post().
//		Resource("pods").
//		Name(pod.Name).
//		Namespace(pod.Namespace).
//		SubResource("exec")
//	req.VersionedParams(&core.PodExecOptions{
//		Container: pod.Spec.Containers[0].Name,
//		Command:   command,
//		Stdout:    true,
//		Stderr:    true,
//	}, scheme.ParameterCodec)
//
//	exec, err := remotecommand.NewSPDYExecutor(f.ClientConfig, "POST", req.URL())
//	if err != nil {
//		return "", fmt.Errorf("failed to init executor: %v", err)
//	}
//
//	err = exec.Stream(remotecommand.StreamOptions{
//		Stdout: &execOut,
//		Stderr: &execErr,
//	})
//
//	if err != nil {
//		return "", fmt.Errorf("could not execute: %v", err)
//	}
//
//	if execErr.Len() > 0 {
//		return "", fmt.Errorf("stderr: %v", execErr.String())
//	}
//
//	return execOut.String(), nil
//}
