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

func initDebugPrint() {
	dPrint = &dPrintItem{
		NewBitmapText("", 0.02, defaultFont, NewColor(100, 100, 100, 255),
			NewScreenPosition(0.01, 0.01, ScreenRelativeLeft)),
		-1,
	}
	dPrintQueue = make([]*dPrintItem, 0, dPrintQueueSize)
}

var appName string = "excavation"

const (
	defaultFont     = "overlays/gui/default/font.material.xml"
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
type dPrintItem struct {
	text  *BitmapText
	timer float64
}

var dPrint *dPrintItem
var dPrintQueue []*dPrintItem
var PrintTime float64 = 5

//Print prints a message to the upper right hand side of the screen.
//  Message will fade after a few seconds,  Any new calls to print will replace
// the current message and reset the timer
// if you want to print a list of messages, use Println
func Print(a ...interface{}) {
	dPrint.text.Text = fmt.Sprint(a...)
	dPrint.timer = Time() + PrintTime
}

//Printf is similar to Print, but accepts a format string
// see http://golang.org/pkg/fmt/#pkg-overview for format details
func Printf(format string, a ...interface{}) {
	dPrint.text.Text = fmt.Sprintf(format, a...)
	dPrint.timer = Time() + PrintTime
}

//Println prints a message to the upper right hand side of the screen
// which will fade after a few seconds or until a new messages pushes it
// off the screen.  Subsequent calls push the previous call down a line
// and show up a the top
func Println(a ...interface{}) {
	dPrintAddToQueue(fmt.Sprint(a...))
}

func dPrintAddToQueue(text string) {
	newItem := &dPrintItem{
		NewBitmapText(text, dPrint.text.Size, dPrint.text.FontMaterial.Name(),
			dPrint.text.Color, NewScreenPosition(0.01, 0.01, ScreenRelativeLeft)),
		Time() + PrintTime,
	}
	dPrintQueue = append(dPrintQueue, newItem)

	if len(dPrintQueue) > dPrintQueueSize {
		over := len(dPrintQueue) - dPrintQueueSize
		dPrintQueue = dPrintQueue[over:]
	}

	for i := range dPrintQueue {
		dPrintQueue[i].text.Position.Y += (dPrintQueue[i].text.Size)
	}

}

//Printfln prints a stack of messages to the screen similar to Println
// but utilizes the format strings of Printf
func Printfln(format string, a ...interface{}) {
	dPrintAddToQueue(fmt.Sprintf(format, a...))
}

func updateDebugPrint() {
	if dPrint.timer >= Time() {
		dPrint.text.Place()
	}

	var timedOut int
	for i := range dPrintQueue {
		if dPrintQueue[i].timer >= Time() {
			dPrintQueue[i].text.Place()
		} else {
			//last index of timed out item
			timedOut = i
		}

	}
	dPrintQueue = dPrintQueue[timedOut:]
}
