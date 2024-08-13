package controllers

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daiki-kim/tweet-app/backend/configs"
	utils "github.com/daiki-kim/tweet-app/backend/pkg"
)

func GoogleLogin(ctx *gin.Context) {
	state, err := utils.GenerateRandomString(32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}
	log.Printf("generate state: %s", state)

	// TODO: set state to cookie or session 2024-08-11
	ctx.SetCookie("oath_state", state, 3600, "/", "", false, true)

	// set redirect url with state
	url := configs.Config.GoogleLoginConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusFound, url)
}

func GoogleCallback(ctx *gin.Context) {
	// get state from cookie
	cookieState, err := ctx.Cookie("oath_state")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get state from cookie"})
		return
	}

	// check state is matched with state from cookie
	state := ctx.Query("state")
	log.Printf("get state: %s", state)

	if state != cookieState {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "state does not match"})
		return
	}

	// get code
	code := ctx.Query("code")
	googleConfig := configs.LoadAppConfig()

	// convert code to token
	token, err := googleConfig.Exchange(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert code to token"})
		return
	}

	// get user info from google api using token
	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	// parse user info
	userData, err := io.ReadAll(res.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse user info"})
	}

	ctx.JSON(http.StatusOK, gin.H{"data": string(userData)})
}
