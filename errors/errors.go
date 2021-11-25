package errors

import (
	"fmt"
	"strconv"

	loghelper "github.com/rayyone/go-core/helpers/log"
	"github.com/rayyone/go-core/helpers/method"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

// ErrorType is the type of an error
type ErrorType uint

// NoType error
const (
	NoType              ErrorType = 500
	InternalServer      ErrorType = 500
	BadRequest          ErrorType = 400
	Unauthorized        ErrorType = 401
	NotFound            ErrorType = 404
	Validation          ErrorType = 422
	UnprocessableEntity ErrorType = 422
	TooManyRequests     ErrorType = 429
)

type CustomError struct {
	errorType     ErrorType
	originalError error
	contexts      []errorContext
	stackTrace    []string
}

type errorContext struct {
	Field   string
	Message string
}

// New creates a new CustomError
func (errorType ErrorType) New(msg string) error {
	loghelper.PrintRed(msg)
	shouldReport := shouldReport(errorType)
	customErr := CustomError{errorType: errorType, originalError: errors.New(msg), stackTrace: []string{msg}}
	if shouldReport {
		customErr.Report()
	}

	return customErr
}

// NewAndReport creates a new CustomError and report
func (errorType ErrorType) NewAndReport(msg string) error {
	loghelper.PrintRed(msg)
	customErr := CustomError{errorType: errorType, originalError: errors.New(msg), stackTrace: []string{msg}}
	customErr.Report()

	return customErr
}

// Newf creates a new CustomError with formatted message
func (errorType ErrorType) Newf(msg string, args ...interface{}) error {
	loghelper.PrintRed(fmt.Sprintf(msg, args...))
	shouldReport := shouldReport(errorType)
	customErr := CustomError{errorType: errorType, originalError: fmt.Errorf(msg, args...), stackTrace: []string{msg}}
	if shouldReport {
		customErr.Report()
	}

	return customErr
}

// NewfAndReport creates a new CustomError with formatted message and report
func (errorType ErrorType) NewfAndReport(msg string, args ...interface{}) error {
	loghelper.PrintRed(fmt.Sprintf(msg, args...))
	customErr := CustomError{errorType: errorType, originalError: fmt.Errorf(msg, args...), stackTrace: []string{msg}}
	customErr.Report()

	return customErr
}

func shouldReport(errorType ErrorType) bool {
	statusCodeStr := strconv.Itoa(int(errorType))[:3] // Get first 3 digits
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		return true
	}

	switch ErrorType(statusCode) {
	case Unauthorized, NotFound, UnprocessableEntity, TooManyRequests:
		return false
	default:
		return true
	}
}

// Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(err error, msg string) error {
	return errorType.Wrapf(err, msg)
}

// Wrapf creates a new wrapped error with formatted message
func (errorType ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	return CustomError{errorType: errorType, originalError: errors.Wrapf(err, msg, args...)}
}

// Error returns the mssage of a CustomError
func (error CustomError) Error() string {
	return error.originalError.Error()
}

func (error CustomError) Report() {
	loghelper.PrintRed("========== Error Stack Strace ==========")
	var stackTrace []string
	for i := 4; i < 9; i++ { // Skip 4 function, Get last 5 error trace
		file, line, fnName := method.TraceCaller(i)
		traceMsg := fmt.Sprintf("%s:%d@%s", file, line, fnName)
		loghelper.PrintYellow(traceMsg)
		stackTrace = append(stackTrace, traceMsg)
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtra("stack_trace", stackTrace)
	})
	sentry.CaptureException(error)
}

// New creates a no type error and report to sentry
func New(msg string) error {
	loghelper.PrintRed(msg)
	err := CustomError{errorType: NoType, originalError: errors.New(msg)}

	err.Report()

	return err
}

// NewAndDontReport creates a new CustomError and don't report it
func NewAndDontReport(msg string) error {
	loghelper.PrintRed(msg)
	err := CustomError{errorType: NoType, originalError: errors.New(msg)}

	return err
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	loghelper.PrintRed(fmt.Sprintf(msg, args...))
	err := CustomError{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}

	err.Report()

	return err
}

func Msg(err error, msg string) error {
	fileName, line, fnName := method.TraceCaller(3)
	errorMsg := fmt.Sprintf("%s:%d@%s()", fileName, line, fnName)
	if customErr, ok := err.(CustomError); ok {
		return CustomError{
			errorType:     customErr.errorType,
			originalError: errors.New(msg),
			contexts:      customErr.contexts,
			stackTrace:    append([]string{errorMsg}, customErr.stackTrace...),
		}
	}

	return CustomError{errorType: NoType, originalError: errors.New(msg), stackTrace: []string{errorMsg}}
}

// Wrap an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Wrapf an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(CustomError); ok {
		return CustomError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			contexts:      customErr.contexts,
			stackTrace:    customErr.stackTrace,
		}
	}

	return CustomError{errorType: NoType, originalError: wrappedError}
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// AddStackTrace an error with format string
func AddStackTrace(err error, msg string) error {
	if customErr, ok := err.(CustomError); ok {
		stackTrace := append([]string{msg}, customErr.stackTrace...)
		return CustomError{errorType: customErr.errorType, originalError: customErr.originalError, contexts: customErr.contexts, stackTrace: stackTrace}
	}

	stackTrace := []string{msg}
	return CustomError{errorType: NoType, originalError: err, stackTrace: stackTrace}
}

// GetStackTrace returns the error stack trace
func GetStackTrace(err error) []string {
	if customErr, ok := err.(CustomError); ok {
		return customErr.stackTrace
	}
	return []string{}
}

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field string, message string) error {
	context := errorContext{Field: field, Message: message}
	if customErr, ok := err.(CustomError); ok {
		contexts := append(customErr.contexts, context)
		return CustomError{errorType: customErr.errorType, originalError: customErr.originalError, contexts: contexts, stackTrace: customErr.stackTrace}
	}

	contexts := []errorContext{context}
	return CustomError{errorType: NoType, originalError: err, contexts: contexts}
}

// GetErrorContexts returns the error context
func GetErrorContexts(err error) map[string]string {
	res := make(map[string]string)
	if customErr, ok := err.(CustomError); ok {
		for _, context := range customErr.contexts {
			res[context.Field] = context.Message
		}
	}
	return res
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(CustomError); ok {
		return customErr.errorType
	}

	return NoType
}

// Is Check if error is the specified error type
func Is(err error, errorType ErrorType) bool {
	if err == nil {
		return false
	}

	errType := NoType
	if customErr, ok := err.(CustomError); ok {
		errType = customErr.errorType
	}

	return errType == errorType
}

// IsNotFound Check if error is NotFound error type
func IsNotFound(err error) bool {
	return Is(err, NotFound)
}

func SetExtra(key string, value interface{}) {
	//@TODO: Need to export to an interface!
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtra(key, value)
	})
}
