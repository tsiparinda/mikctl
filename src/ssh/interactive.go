package ssh

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// RunSSHInteractive runs an interactive command on RouterOS over SSH.
// input is sent to the remote command's stdin.
func RunSSHInteractive(ip string, command string, input string) (string, error) {

	cmd := exec.Command(
		"ssh",
		"-tt",
		"-o", "ConnectTimeout=10",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "LogLevel=ERROR",
		ip,
		command,
	)

	var output strings.Builder

	cmd.Stdout = &output
	cmd.Stderr = &output

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("ssh start: %w", err)
	}

	// Give SSH/MAC-Telnet time to start and display Login:
	time.Sleep(500 * time.Millisecond)

	if _, err := stdin.Write([]byte(input)); err != nil {
		return output.String(), fmt.Errorf("stdin write: %w", err)
	}

	stdin.Close()

	err = cmd.Wait()

	return strings.TrimSpace(output.String()), err
}