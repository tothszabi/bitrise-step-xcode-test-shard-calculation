package step

import (
	"fmt"
	"os"
	"os/exec"
)

func CreateTempFolder() (string, error) {
	return os.MkdirTemp("", "")
}

func Execute(cmd *exec.Cmd) error {
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return err
	}

	return nil
}
