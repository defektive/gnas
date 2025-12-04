package process

import (
	"os/user"
	"strconv"
)

func IsAdmin() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}

	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		return false
	}

	return uid == 0
}
