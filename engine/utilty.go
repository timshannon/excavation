package engine

import (
	"os"
	"os/user"
	"path"
)

//UserDir is the current users folder where everything user specific
// like controls, save games, and video settings will be stored
// if the path isn't found, it'll be created
func UserDir() (string, error) {
	var userDir string
	curUser, err := user.Current()
	if err != nil {
		return "", err
	}
	userDir = path.Join(curUser.HomeDir, ".excavation")
	if err := os.Chdir(userDir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(userDir, 0774)
		} else {
			return "", err
		}
	}

	return userDir, nil
}
