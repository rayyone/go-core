package response

import (
	"net/http"
	"strconv"

	"github.com/rayyone/go-core/errors"
	"github.com/rayyone/go-core/helpers/pagination"
	"github.com/gin-gonic/gin"
)

type StandardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResponseWithData struct {
	StandardResponse
	Data interface{} `json:"data"`
}

type PaginatorResponse struct {
	ResponseWithData
	Paginator *pagination.Paginator `json:"paginator"`
}

type ErrorResponse struct {
	StandardResponse
	ErrorCode string            `json:"error_code"`
	Errors    []string          `json:"errors"`
	Fields    map[string]string `json:"fields"`
}

func BuildStandardResponse(status string, message string) StandardResponse {
	var response StandardResponse
	response.Status = status
	response.Message = message

	return response
}

func BuildSuccessResponse(message string, data interface{}) ResponseWithData {
	var response ResponseWithData
	response.StandardResponse = BuildStandardResponse("success", message)
	response.Data = data

	return response
}

func BuildErrorResponse(err error, errorCode string, message string) ErrorResponse {
	var response ErrorResponse
	var errMsgs []string

	response.StandardResponse = BuildStandardResponse("error", message)
	contexts := errors.GetErrorContexts(err)
	for _, v := range contexts {
		errMsgs = append(errMsgs, v)
	}
	if response.Errors = errMsgs; errMsgs == nil {
		response.Errors = []string{message}
	}
	response.Fields = contexts
	response.ErrorCode = errorCode

	return response
}

// RespondSuccess respond JSON with data
func RespondSuccess(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, BuildSuccessResponse(message, data))
}

// RespondSuccessWithPaginator respond with paginator
func RespondSuccessWithPaginator(c *gin.Context, data interface{}, paginator *pagination.Paginator, message string) {
	var response PaginatorResponse
	response.StandardResponse = BuildStandardResponse("success", message)
	response.Paginator = paginator
	response.Data = data

	c.JSON(http.StatusOK, response)
}

// RespondError respond error
func RespondError(c *gin.Context, err error) {
	var defaultMessage, errorCode string

	errType := errors.GetType(err)
	switch errType {
	case errors.Unauthorized:
		defaultMessage = "Unauthorized."
	case errors.NotFound:
		defaultMessage = "Resource not found."
	case errors.UnprocessableEntity:
		defaultMessage = "Unprocessable entity error."
	case errors.BadRequest:
		defaultMessage = "Bad request."
	default:
		defaultMessage = "Internal server error."
	}

	message := err.Error()
	if message == "" {
		message = defaultMessage
	}
	if errorCode == "" {
		// Get error code from err type value
		errorCode = strconv.Itoa(int(errType))
	}
	statusCode := getStatusCode(errorCode)
	c.JSON(statusCode, BuildErrorResponse(err, errorCode, message))
}

func getStatusCode(errorCode string) (statusCode int) {
	if len(errorCode) == 3 {
		statusCode, _ = strconv.Atoi(errorCode)
	} else if len(errorCode) > 3 {
		errorCodeStr := errorCode[:3]
		statusCode, _ = strconv.Atoi(errorCodeStr)
	}
	if statusCode == 0 {
		statusCode = 500
	}

	return statusCode
}
