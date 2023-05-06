// Package docker provides support for starting and stopping docker containers
// for running tests.
package docker

import (
	"bytes"
	"encoding/json"
	"net"
	"os/exec"
	"testing"
)

// Container tracks information about the docker container for tests.
type Container struct {
	ID   string
	Host string // IP:Port
}

// StartContainer starts the specified container for running tests.
func StartContainer(t *testing.T, image string, port string, args ...string) *Container {
	arg := []string{"run", "-P", "-d"}
	arg = append(arg, args...)
	arg = append(arg, image)

	cmd := exec.Command("docker", arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not start container %s: %v", image, err)
	}

	id := out.String()[:12]

	cmd = exec.Command("docker", "inspect", id)
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not inspect container %s: %v", id, err)
	}

	var doc []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("could not decode json: %v", err)
	}

	ip, randPort := extractIPPort(doc, port)

	c := Container{
		ID:   id,
		Host: net.JoinHostPort(ip, randPort),
	}

	t.Logf("Image: %s", image)
	t.Logf("ContainerID: %s", c.ID)
	t.Logf("Host: %s", c.Host)

	return &c

}

// StopContainer stops and removes the specified container.
func StopContainer(t *testing.T, id string) {
	if err := exec.Command("docker", "stop", id).Run(); err != nil {
		t.Fatalf("could not stop container: %v", err)
	}
	t.Log("Stopped:", id)

	if err := exec.Command("docker", "rm", id, "-v").Run(); err != nil {
		t.Fatalf("coould not remove container: %v", err)
	}
	t.Log("Removed:", id)
}

// DumpContainerLogs outputs logs from the running docker container.
func DumpContainerLogs(t *testing.T, id string) {
	out, err := exec.Command("docker", "logs", id).CombinedOutput()
	if err != nil {
		t.Fatalf("could not log container: %v", err)
	}
	t.Logf("Logss for %s\n%s", id, out)
}

func extractIPPort(doc []map[string]interface{}, port string) (string, string) {
	nw, exists := doc[0]["NetworkSettings"]
	if !exists {
		return "", ""
	}
	ports, exists := nw.(map[string]interface{})["Ports"]
	if !exists {
		return "", ""
	}
	tcp, exists := ports.(map[string]interface{})[port+"/tcp"]
	if !exists {
		return "", ""
	}
	list, exists := tcp.([]interface{})
	if !exists {
		return "", ""
	}

	var hostIP string
	var hostPort string
	for _, l := range list {
		data, exists := l.(map[string]interface{})
		if !exists {
			return "", ""
		}
		hostIP = data["HostIp"].(string)
		if hostIP != "::" {
			hostPort = data["HostPort"].(string)
		}
	}

	return hostIP, hostPort
}
