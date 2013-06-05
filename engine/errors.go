// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package engine

//This is an error handler which collects all of the errors
// during the engine processing and lets you handle
// them however you want.  Write to log file
// print to console, etc
// Engine stopping errors (like during initialization) will panic
// the rest will be added here and be returned in their
// respective functions

var errorList []error

var errorHandler ErrorHandler

type ErrorHandler func(e error)

func SetErrorHandler(f ErrorHandler) {
	errorHandler = f
}

func RaiseError(e error) {
	if errorHandler != nil {
		errorHandler(e)
	}

	errorList = append(errorList, e)
}

func GetErrors() []error {
	return errorList
}
