package form3

import (
	"os/exec"
	"strings"
)

const (
	uuidChars = "abcdefghijkmnpqrstuvwxyz0123456789"
)

// todo: replace with github.com/google/uuid
func CreateUUID() string {
	out, _ := exec.Command("uuidgen").Output()
	return strings.TrimRight(string(out), "\n")
}
