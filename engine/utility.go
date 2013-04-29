// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package engine

import (
	"fmt"
	"os"
	"os/user"
	"path"
)

var appName string = "excavation"

const (
	defaultFont     = "fonts/ubuntu/Ubuntu-R.ttf"
	dPrintQueueSize = 30
)

//UserDir is the current users folder where everything user specific
// like controls, save games, and video settings will be stored
// if the path isn't found, it'll be created
func UserDir() (string, error) {
	var userDir string

	//TODO: Testing and specific handling on Windows and Mac
	userDir = os.Getenv("XDG_DATA_HOME")
	if userDir == "" {
		curUser, err := user.Current()
		if err != nil {
			return "", err
		}
		userDir = path.Join(curUser.HomeDir, ".local/share/"+appName)
	} else {
		userDir = path.Join(userDir, appName)
	}

	if err := os.Chdir(userDir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(userDir, 0774)
		} else {
			return "", err
		}
	}

	return userDir, nil
}

//Debug Printing
var dPrintText *Text
var dPrintString []string
var dPrintTimer []float64
var PrintTime float64 = 120

//Print prints a message to the upper right hand side of the screen.
//  Message will fade after a few seconds,  Any new calls to print will replace
// the current message and reset the timer
// if you want to print a list of messages, use Println
func Print(a ...interface{}) {
	initDebugPrint()
	dPrintString[0] = fmt.Sprint(a...)
	dPrintTimer[0] = Time() + PrintTime
	dPrintText.SetText(dPrintString)
}

func initDebugPrint() {
	if dPrintText == nil || !dPrintText.Overlay().Material.IsValid() {
		dPrintString = make([]string, dPrintQueueSize)
		dPrintTimer = make([]float64, dPrintQueueSize)
		dPrintText = NewText(dPrintString, defaultFont, 18, NewColor(75, 75, 75, 255),
			NewScreenArea(0, 0, 1, 1, ScreenRelativeLeft))
	}
}

//Printf is similar to Print, but accepts a format string
// see http://golang.org/pkg/fmt/#pkg-overview for format details
func Printf(format string, a ...interface{}) {
	initDebugPrint()
	dPrintString[0] = fmt.Sprintf(format, a...)
	dPrintTimer[0] = Time() + PrintTime
	dPrintText.SetText(dPrintString)
}

//Println prints a message to the upper right hand side of the screen
// which will fade after a few seconds or until a new messages pushes it
// off the screen.  Subsequent calls push the previous call down a line
// and show up a the top
func Println(a ...interface{}) {
	dPrintAddToQueue(fmt.Sprint(a...))
}

func dPrintAddToQueue(text string) {
	initDebugPrint()

	dPrintString = append(dPrintString, text)
	dPrintTimer = append(dPrintTimer, Time()+PrintTime)

	if len(dPrintString) > dPrintQueueSize {
		over := len(dPrintString) - dPrintQueueSize
		dPrintString = dPrintString[over:]
		dPrintTimer = dPrintTimer[over:]
	}

	dPrintText.SetText(dPrintString)
}

//Printfln prints a stack of messages to the screen similar to Println
// but utilizes the format strings of Printf
func Printfln(format string, a ...interface{}) {
	dPrintAddToQueue(fmt.Sprintf(format, a...))
}

func updateDebugPrint() {
	if dPrintText != nil && len(dPrintString) > 0 {
		var timedOut int
		for i := range dPrintString {
			if dPrintTimer[i] < Time() {
				//last index of timed out item
				timedOut = i
			}

		}
		if timedOut != 0 {
			dPrintTimer = dPrintTimer[timedOut:]
			dPrintString = dPrintString[timedOut:]
			dPrintText.SetText(dPrintString)
		}
		dPrintText.Place()
	}
}
