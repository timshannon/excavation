// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

package engine

import (
	"os"
	"os/user"
	"path"
)

var appName string = "excavation"

//UserDir is the current users folder where everything user specific
// like controls, save games, and video settings will be stored
// if the path isn't found, it'll be created
func UserDir() (string, error) {
	var userDir string
	curUser, err := user.Current()
	if err != nil {
		return "", err
	}
	userDir = path.Join(curUser.HomeDir, "."+appName)
	if err := os.Chdir(userDir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(userDir, 0774)
		} else {
			return "", err
		}
	}

	return userDir, nil
}
