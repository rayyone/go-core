package middleware

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
	"unicode"

	"github.com/rayyone/go-core/errors"
	"github.com/rayyone/go-core/helpers/response"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func isAllUpper(s string) bool {
	for _, v := range s {
		if !unicode.IsUpper(v) {
			return false
		}
	}
	return true
}

func lowerCaseFirst(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	for i, v := range s {
		return string(unicode.ToLower(v)) + s[i+1:]
	}
	return ""
}

func fieldErrorToText(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than or equal to %s", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than or equal to %s", e.Field(), e.Param())
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field(), e.Param())
	}
	return fmt.Sprintf("%s is not valid", e.Field())
}

func handleValidationError(e *gin.Error, c *gin.Context) {
	validationErrs := e.Err.(validator.ValidationErrors)
	var err error
	var errMessage string
	for _, validationErr := range validationErrs {
		errMessage = fieldErrorToText(validationErr)
		if isAllUpper(validationErr.Field()) {
			err = errors.AddErrorContext(err, strings.ToLower(validationErr.Field()), errMessage)
		} else {
			err = errors.AddErrorContext(err, lowerCaseFirst(validationErr.Field()), errMessage)
		}
	}
	response.RespondError(c, errors.Validation.New(errMessage))
}

// HandleError Middleware for handling error
func HandleError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, ginErr := range c.Errors {
			switch e := ginErr.Err.(type) {
			case validator.ValidationErrors:
				handleValidationError(ginErr, c)
			case *net.OpError:
				if se, ok := e.Err.(*os.SyscallError); ok {
					if se.Err == syscall.EPIPE {
						response.RespondError(c, errors.NewAndDontReport("Broken Pipe"))
						log.Printf("Error: Broken Pipe | %+v", ginErr)
					} else if se.Err == syscall.ECONNRESET {
						response.RespondError(c, errors.NewAndDontReport("Connection Reset"))
						log.Printf("Error: Connection Reset | %+v", ginErr)
					}
				}
			default:
				err := errors.Newf("Unknown error. Error: %s", spew.Sdump(e))
				response.RespondError(c, errors.Msg(err, "Unknown error"))
				log.Printf("Error: Unknown error | %v", ginErr)
			}
		}

		// If there is no response yet, we respond with unhandled response error
		if !c.Writer.Written() {
			response.RespondError(c, errors.New("Unhandled response."))
			log.Println("Error: Unhandled response.")
		}
	}
}
