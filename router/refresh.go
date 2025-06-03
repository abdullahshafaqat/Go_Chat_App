package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *routerImpl) RefreshToken(c *gin.Context) {
	refreshToken := c.PostForm("refresh_token")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}

	ID, err := r.authservice.VerifyRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	newAccessToken, _, err := r.authservice.GenerateTokens(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})

}
