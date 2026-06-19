package ssh

import (
    "fmt"
    "os/exec"
)

func CopyFile(
    ip string,
    localFile string,
    remoteFile string,
) error {

    out, err := exec.Command(
        "scp",
        localFile,
        fmt.Sprintf("%s:%s", ip, remoteFile),
    ).CombinedOutput()

    if err != nil {
        return fmt.Errorf(
            "scp failed: %s: %w",
            string(out),
            err,
        )
    }

    return nil
}