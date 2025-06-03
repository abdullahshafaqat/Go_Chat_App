package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *routerImpl) Authorize(c *gin.Context) {
	token := r.authservice.BearerToken(c.GetHeader("Authorization"))

	valid, msg, err := r.authservice.Authorize(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": msg, 
		})
		return
	}
	c.Status(http.StatusOK)
}
