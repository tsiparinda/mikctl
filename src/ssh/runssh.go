package ssh

import (
	"fmt"
	"os/exec"
	"strings"
)

func RunSSH(ip string, command string) (string, error) {

	out, err := exec.Command(
		"ssh",
		"-o", "ConnectTimeout=10",
		"-o", "BatchMode=yes",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "LogLevel=ERROR",
		ip,
		command,
	).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), err
}
