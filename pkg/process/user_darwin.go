package process

import (
	"os/exec"
	"strings"
)

func IsAdmin() bool {
	cmd := exec.Command("groups")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	groups := strings.Fields(string(output))
	for _, group := range groups {
		if group == "admin" {
			return true
		}
	}
	return false
}
