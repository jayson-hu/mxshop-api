package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "OPTIONS,DELETE,POST,GET,PUT,PATCH")
		c.Header("Access-Control-Expose-Headers", "Content-Disposition,Message,Code,Access-Token")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "X-Custom-Header,accept,Content-Type,Access-Token,Message,Code,Token,X-Token")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

	}

}
