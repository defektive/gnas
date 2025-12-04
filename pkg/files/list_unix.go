//go:build linux || darwin

package files

import (
	"os"
	"syscall"
)

// getPermissionsString creates a Unix-like permission string from a file mode.
func getPermissionsString(mode os.FileMode) string {
	perms := ""
	if mode&syscall.S_IRUSR != 0 {
		perms += "r"
	} else {
		perms += "-"
	}
	if mode&syscall.S_IWUSR != 0 {
		perms += "w"
	} else {
		perms += "-"
	}
	if mode&syscall.S_IXUSR != 0 {
		perms += "x"
	} else {
		perms += "-"
	}

	if mode&syscall.S_IRGRP != 0 {
		perms += "r"
	} else {
		perms += "-"
	}
	if mode&syscall.S_IWGRP != 0 {
		perms += "w"
	} else {
		perms += "-"
	}
	if mode&syscall.S_IXGRP != 0 {
		perms += "x"
	} else {
		perms += "-"
	}

	if mode&syscall.S_IROTH != 0 {
		perms += "r"
	} else {
		perms += "-"
	}
	if mode&syscall.S_IWOTH != 0 {
		perms += "w"
	} else {
		perms += "-"
	}
	if mode&syscall.S_IXOTH != 0 {
		perms += "x"
	} else {
		perms += "-"
	}
	return perms
}
