package delivery

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error string `json:"message,omitempty"`
}

func SetError(c *gin.Context, statusCode int, errorMsg string) {
	c.JSON(statusCode, ErrorResponse{Error: errorMsg})
}
