package ryerr

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	loghelper "github.com/rayyone/go-core/helpers/log"
	ry_slack "github.com/rayyone/go-core/helpers/slack"
	"gorm.io/gorm"
)

// ErrorType is the type of an error
type ErrorType uint

// NoType error
const (
	NoType              ErrorType = 500
	InternalServer      ErrorType = 500
	BadRequest          ErrorType = 400
	Unauthorized        ErrorType = 401
	Forbidden           ErrorType = 403
	NotFound            ErrorType = 404
	Validation          ErrorType = 422
	UnprocessableEntity ErrorType = 422
	TooManyRequests     ErrorType = 429
)

type Err struct {
	errorType     ErrorType
	originalError error
	contexts      []errorContext
	stackTrace    []string
}

type errorContext struct {
	Field   string
	Message string
}

func (c Err) Error() string {
	return c.originalError.Error()
}

// New creates a new Err
func (errorType ErrorType) New(msg string) error {
	loghelper.PrintRed(msg)
	shouldReport := shouldReport(errorType)
	errLogger := gin.DefaultErrorWriter
	errLogger.Write([]byte(msg + "\n"))

	customErr := Err{errorType: errorType, originalError: errors.New(msg), stackTrace: []string{msg}}
	if shouldReport {
		customErr.Report()
	}

	return customErr
}

// NewAndReport creates a new Err and report
func (errorType ErrorType) NewAndReport(msg string) error {
	loghelper.PrintRed(msg)
	errLogger := gin.DefaultErrorWriter
	errLogger.Write([]byte(msg + "\n"))
	customErr := Err{errorType: errorType, originalError: errors.New(msg), stackTrace: []string{msg}}
	customErr.Report()

	return customErr
}

// Newf creates a new Err with formatted message
func (errorType ErrorType) Newf(msg string, args ...interface{}) error {
	out := fmt.Sprintf(msg, args...)
	loghelper.PrintRed(out)
	shouldReport := shouldReport(errorType)
	errLogger := gin.DefaultErrorWriter
	errLogger.Write([]byte(out + "\n"))

	customErr := Err{errorType: errorType, originalError: fmt.Errorf(msg, args...), stackTrace: []string{msg}}
	if shouldReport {
		customErr.Report()
	}

	return customErr
}

// NewfAndReport creates a new Err with formatted message and report
func (errorType ErrorType) NewfAndReport(msg string, args ...interface{}) error {
	out := fmt.Sprintf(msg, args...)
	loghelper.PrintRed(out)
	errLogger := gin.DefaultErrorWriter
	errLogger.Write([]byte(out + "\n"))

	customErr := Err{errorType: errorType, originalError: fmt.Errorf(msg, args...), stackTrace: []string{msg}}
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
	return Err{errorType: errorType, originalError: errors.Wrapf(err, msg, args...)}
}

func (c Err) Report() {
	var fullText = "========== Error Stack Strace =========="
	loghelper.PrintRed(fullText)
	fullText = "\n" + fullText + "\n"
	var stackTrace []string
	for i := 4; i < 9; i++ { // Skip 4 function, Get last 5 error trace
		file, line, fnName := traceCaller(i)
		traceMsg := fmt.Sprintf("%s:%d@%s", file, line, fnName)
		loghelper.PrintYellow(traceMsg)
		stackTrace = append(stackTrace, traceMsg)
		fullText += traceMsg + "\n"
	}
	errLogger := gin.DefaultErrorWriter
	errLogger.Write([]byte(fullText))

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtra("stack_trace", stackTrace)

	})
	eventId := sentry.CaptureException(c)
	slackMsg := fmt.Sprintf("*%s*\n", c.Error())
	if eventId != nil {
		slackMsg += fmt.Sprintf("*EventID:* %s\n", *eventId)
		if sentryProjectUrl := ry_slack.CurrentSlackClient().GetOption("sentry_project_url"); sentryProjectUrl != nil {
			slackMsg += fmt.Sprintf("<%s?query=%s|See more detail>", sentryProjectUrl, *eventId)
		}
	}
	if errChannel := ry_slack.CurrentSlackClient().GetOption("error_channel"); errChannel != nil {
		ry_slack.SendSimpleMessageToChannel(errChannel.(string), "ry-api error", slackMsg)
	} else {
		ry_slack.SendSimpleMessage("ry-api error", slackMsg)
	}
}

// New creates a no type error and report to sentry
func New(msg string) error {
	loghelper.PrintRed(msg)
	errLogger := gin.DefaultErrorWriter
	errLogger.Write([]byte(msg + "\n"))
	err := Err{errorType: NoType, originalError: errors.New(msg)}

	err.Report()

	return err
}

// NewAndDontReport creates a new Err and don't report it
func NewAndDontReport(msg string) error {
	loghelper.PrintRed(msg)
	errLogger := gin.DefaultErrorWriter
	errLogger.Write([]byte(msg + "\n"))
	err := Err{errorType: NoType, originalError: errors.New(msg)}

	return err
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	out := fmt.Sprintf(msg, args...)
	loghelper.PrintRed(out)
	errLogger := gin.DefaultErrorWriter
	errLogger.Write([]byte(out + "\n"))
	err := Err{errorType: NoType, originalError: errors.New(out)}

	err.Report()

	return err
}

func Msg(err error, msg string) error {
	fileName, line, fnName := traceCaller(3)
	errorMsg := fmt.Sprintf("%s:%d@%s()", fileName, line, fnName)
	if customErr, ok := err.(Err); ok {
		return Err{
			errorType:     customErr.errorType,
			originalError: errors.New(msg),
			contexts:      customErr.contexts,
			stackTrace:    append([]string{errorMsg}, customErr.stackTrace...),
		}
	}

	return Err{errorType: NoType, originalError: errors.New(msg), stackTrace: []string{errorMsg}}
}

// Wrap an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Wrapf an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(Err); ok {
		return Err{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			contexts:      customErr.contexts,
			stackTrace:    customErr.stackTrace,
		}
	}

	return Err{errorType: NoType, originalError: wrappedError}
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err).(Err)
}

// AddStackTrace an error with format string
func AddStackTrace(err error, msg string) error {
	if customErr, ok := err.(Err); ok {
		stackTrace := append([]string{msg}, customErr.stackTrace...)
		return Err{errorType: customErr.errorType, originalError: customErr.originalError, contexts: customErr.contexts, stackTrace: stackTrace}
	}

	stackTrace := []string{msg}
	return Err{errorType: NoType, originalError: err, stackTrace: stackTrace}
}

// GetStackTrace returns the error stack trace
func GetStackTrace(err error) []string {
	if customErr, ok := err.(Err); ok {
		return customErr.stackTrace
	}
	return []string{}
}

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field string, message string) error {
	context := errorContext{Field: field, Message: message}
	if customErr, ok := err.(Err); ok {
		contexts := append(customErr.contexts, context)
		return Err{errorType: customErr.errorType, originalError: customErr.originalError, contexts: contexts, stackTrace: customErr.stackTrace}
	}

	contexts := []errorContext{context}
	return Err{errorType: NoType, originalError: err, contexts: contexts}
}

// GetErrorContexts returns the error context
func GetErrorContexts(err error) map[string]string {
	res := make(map[string]string)
	if customErr, ok := err.(Err); ok {
		for _, context := range customErr.contexts {
			res[context.Field] = context.Message
		}
	}
	return res
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(Err); ok {
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
	if customErr, ok := err.(Err); ok {
		errType = customErr.errorType
	}

	return errType == errorType
}

// IsNotFound Check if error is NotFound error type
func IsNotFound(err error) bool {
	return Is(err, NotFound)
}

func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IsGormError(err error) bool {
	return err != nil && !IsRecordNotFound(err)
}

func SetExtra(key string, value interface{}) {
	//@TODO: Need to export to an interface!
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtra(key, value)
	})
}

func traceCaller(skip int) (file string, line int, fnName string) {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(skip, pc)

	if pc[0] == uintptr(0) {
		return
	}

	f := runtime.FuncForPC(pc[0])
	file, line = f.FileLine(pc[0])
	fnName = f.Name()

	return
}
