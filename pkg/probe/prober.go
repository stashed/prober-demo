/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package probe

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"stash.appscode.dev/prober-demo/api/v1"
	"strconv"
	"strings"
	"time"

	"k8s.io/client-go/rest"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/klog"
)

// Result is a string used to handle the results for probing container readiness/liveness
type Result string

const (
	// Success Result
	Success Result = "success"
	// Warning Result. Logically success, but with additional debugging information attached.
	Warning Result = "warning"
	// Failure Result
	Failure Result = "failure"
	// Unknown Result
	Unknown Result = "unknown"
)

type prober struct {
	httpGet HttpGetProber
	httpPost HttpPostProber
	tcp     TcpProber
	exec    ExecProber
	config  *rest.Config
}

// newProber creates a prober instance that can be used to run httpGet, tcp or exec probe.
func newProber(config *rest.Config) *prober {
	const followNonLocalRedirects = false

	return &prober{
		httpGet: NewHTTPGetProber(followNonLocalRedirects),
		httpPost: NewHTTPPostProber(followNonLocalRedirects),
		tcp:     NewTcpProber(),
		exec:    NewExecProber(),
		config:  config,
	}
}

// buildHeaderMap takes a list of HTTPHeader <name, value> string
// pairs and returns a populated string->[]string http.Header map.
func buildHeader(headerList []v1.HTTPHeader) http.Header {
	headers := make(http.Header)
	for _, header := range headerList {
		headers[header.Name] = append(headers[header.Name], header.Value)
	}
	return headers
}

func (pb *prober) runProbe(p *v1.Probe, pod *core.Pod, status core.PodStatus, container core.Container) (Result, string, error) {
	timeout := time.Duration(p.TimeoutSeconds) * time.Second
	if p.Exec != nil {
		klog.V(4).Infof("Exec-Probe Pod: %v, Container: %v, Command: %v", pod, container, p.Exec.Command)
		return pb.exec.Probe(pb.config, pod, container, p.Exec.Command)
	}
	if p.HTTPGet != nil {
		scheme := strings.ToLower(string(p.HTTPGet.Scheme))
		host := p.HTTPGet.Host
		if host == "" {
			host = status.PodIP
		}
		port, err := extractPort(p.HTTPGet.Port, container)
		if err != nil {
			return Unknown, "", err
		}
		path := p.HTTPGet.Path
		klog.V(4).Infof("HTTP-Probe Host: %v://%v, Port: %v, Path: %v", scheme, host, port, path)
		targetURL := formatURL(scheme, host, port, path)
		headers := buildHeader(p.HTTPGet.HTTPHeaders)
		klog.V(4).Infof("HTTP-Probe Headers: %v", headers)
		return pb.httpGet.Probe(targetURL, headers, timeout)
	}
	if p.HTTPPost != nil {
		scheme := strings.ToLower(string(p.HTTPPost.Scheme))
		host := p.HTTPPost.Host
		if host == "" {
			host = status.PodIP
		}
		port, err := extractPort(p.HTTPPost.Port, container)
		if err != nil {
			return Unknown, "", err
		}
		path := p.HTTPPost.Path
		klog.V(4).Infof("HTTP-Probe Host: %v://%v, Port: %v, Path: %v", scheme, host, port, path)
		targetURL := formatURL(scheme, host, port, path)
		headers := buildHeader(p.HTTPPost.HTTPHeaders)
		klog.V(4).Infof("HTTP-Probe Headers: %v", headers)
		return pb.httpPost.Probe(targetURL, headers, p.HTTPPost.Form,p.HTTPPost.Body, timeout)
	}
	if p.TCPSocket != nil {
		port, err := extractPort(p.TCPSocket.Port, container)
		if err != nil {
			return Unknown, "", err
		}
		host := p.TCPSocket.Host
		if host == "" {
			host = status.PodIP
		}
		klog.V(4).Infof("TCP-Probe Host: %v, Port: %v, Timeout: %v", host, port, timeout)
		return pb.tcp.Probe(host, port, timeout)
	}
	klog.Warningf("Failed to find probe builder for container: %v", container)
	return Unknown, "", fmt.Errorf("missing probe handler for %s:%s", Pod(pod), container.Name)
}

func extractPort(param intstr.IntOrString, container core.Container) (int, error) {
	port := -1
	var err error
	switch param.Type {
	case intstr.Int:
		port = param.IntValue()
	case intstr.String:
		if port, err = findPortByName(container, param.StrVal); err != nil {
			// Last ditch effort - maybe it was an int stored as string?
			if port, err = strconv.Atoi(param.StrVal); err != nil {
				return port, err
			}
		}
	default:
		return port, fmt.Errorf("intOrString had no kind: %+v", param)
	}
	if port > 0 && port < 65536 {
		return port, nil
	}
	return port, fmt.Errorf("invalid port number: %v", port)
}

// findPortByName is a helper function to look up a port in a container by name.
func findPortByName(container core.Container, portName string) (int, error) {
	for _, port := range container.Ports {
		if port.Name == portName {
			return int(port.ContainerPort), nil
		}
	}
	return 0, fmt.Errorf("port %s not found", portName)
}

// formatURL formats a URL from args.  For testability.
func formatURL(scheme string, host string, port int, path string) *url.URL {
	u, err := url.Parse(path)
	// Something is busted with the path, but it's too late to reject it. Pass it along as is.
	if err != nil {
		u = &url.URL{
			Path: path,
		}
	}
	u.Scheme = scheme
	u.Host = net.JoinHostPort(host, strconv.Itoa(port))
	return u
}
