package probe

import (
	"fmt"
	"k8s.io/klog"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"strings"
	"time"
)

// ContainerID is a type that identifies a container.
type ContainerID struct {
	// The type of the container runtime. e.g. 'docker'.
	Type string
	// The identification of the container, this is comsumable by
	// the underlying container runtime. (Note that the container
	// runtime interface still takes the whole struct as input).
	ID string
}

// Convenience method for creating a ContainerID from an ID string.
func ParseContainerID(containerID string) ContainerID {
	var id ContainerID
	if err := id.ParseString(containerID); err != nil {
		klog.Error(err)
	}
	return id
}

func (c *ContainerID) ParseString(data string) error {
	// Trim the quotes and split the type and ID.
	parts := strings.Split(strings.Trim(data, "\""), "://")
	if len(parts) != 2 {
		return fmt.Errorf("invalid container ID: %q", data)
	}
	c.Type, c.ID = parts[0], parts[1]
	return nil
}

func (c *ContainerID) String() string {
	return fmt.Sprintf("%s://%s", c.Type, c.ID)
}

func (c *ContainerID) IsEmpty() bool {
	return *c == ContainerID{}
}

func (c *ContainerID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", c.String())), nil
}

func (c *ContainerID) UnmarshalJSON(data []byte) error {
	return c.ParseString(string(data))
}

type ContainerCommandRunner interface {
	// RunInContainer synchronously executes the command in the container, and returns the output.
	// If the command completes with a non-0 exit code, a k8s.io/utils/exec.ExitError will be returned.
	RunInContainer(id kubecontainer.ContainerID, cmd []string, timeout time.Duration) ([]byte, error)
}
