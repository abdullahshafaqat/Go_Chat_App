package router

import (
	"log"
	"net/http"

	"github.com/abdullahshafaqat/Go_Chat_App.git/middelwares"
	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

func (r *routerImpl) Login(c *gin.Context) {
    var input models.UserLogin
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
    }

    userID, err := r.authservice.Login(c, input.Email, input.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    accessToken, refreshToken, err := middelwares.GenerateTokens(userID)
    if err != nil {
        log.Printf("Token generation failed for %s: %v", input.Email, err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
        return
    }

    log.Printf("Generated tokens for %s (ID: %s)", input.Email, userID)

    c.JSON(http.StatusOK, gin.H{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}