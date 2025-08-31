package http

import "github.com/gin-gonic/gin"

func ReturnGinError(c *gin.Context, err error, meta any) {
	c.Errors = append(c.Errors, &gin.Error{
		Err:  err,
		Type: gin.ErrorTypePublic,
		Meta: meta,
	})

	c.Abort()
}
