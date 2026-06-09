package judge

import (
	"fmt"
	"os/exec"
)

const BuildNetwork = "regs-build-net"

func EnsureBuildNetwork() error {
	if err := exec.Command("docker", "network", "inspect", BuildNetwork).Run(); err == nil {
		return nil
	}

	if err := exec.Command("docker", "network", "create", BuildNetwork).Run(); err != nil {
		return fmt.Errorf("無法建立編譯用 Docker 網路 %s: %w", BuildNetwork, err)
	}
	return nil
}
