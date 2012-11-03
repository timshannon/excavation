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

func RegisterErrorHandler(f ErrorHandler) {
	errorHandler = errorHandler
}

func addError(e error) {
	if errorHandler != nil {
		errorHandler(e)
	}

	errorList = append(errorList, e)
}

func GetErrors() []error {
	return errorList
}
