package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daiki-kim/tweet-app/backend/configs"
	"github.com/daiki-kim/tweet-app/backend/pkg/utils"
)

func GoogleLogin(c *gin.Context) {
	state, err := utils.GenerateRandomString(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}

	// TODO: set state to cookie or session 2024-08-11

	url := configs.GoogleConfig().AuthCodeURL(state)
	c.Redirect(http.StatusFound, url)
}
