package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorData struct {
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}

// RespondSuccess - Standard success response
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, SuccessData{
		Message: message,
		Data:    data,
	})
}

// RespondCreated - For 201 created
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, SuccessData{
		Message: message,
		Data:    data,
	})
}

// RespondError - Standard error response
func Error(c *gin.Context, status int, err string, details ...any) {
	var detail any
	if len(details) > 0 {
		detail = details[0]
	}
	c.AbortWithStatusJSON(status, ErrorData{
		Error:   err,
		Details: detail,
	})
}
