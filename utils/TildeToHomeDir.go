// Utilities for the potato system
package utils

import (
	"os/user"
	"strings"
)

func TildeToHomeDir(path string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	if strings.HasPrefix(path, "~") {
		return usr.HomeDir + path[1:]
	}
	return path
}
