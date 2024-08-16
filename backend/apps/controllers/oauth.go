package controllers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/daiki-kim/tweet-app/backend/configs"
	utils "github.com/daiki-kim/tweet-app/backend/pkg"
)

func GoogleLogin(ctx *gin.Context, action string) {
	configs.LoadAppConfig()
	conf := configs.Config

	state, err := utils.GenerateRandomString(32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}

	// add action to state
	stateData := map[string]string{
		"state":  state,
		"action": action,
	}
	stateJSON, err := json.Marshal(stateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal state"})
		return
	}
	encodedState := base64.URLEncoding.EncodeToString(stateJSON)

	ctx.SetCookie("oath_state", encodedState, 3600, "/", "", false, true)

	// set redirect url with state
	url := conf.GoogleLoginConfig.AuthCodeURL(encodedState)
	log.Println(url)
	ctx.Redirect(http.StatusFound, url)
}

func GoogleCallback(ctx *gin.Context) {
	googleConfig := configs.LoadAppConfig()
	conf := configs.Config

	// get state from cookie
	cookieState, err := ctx.Cookie("oath_state")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get state from cookie"})
		return
	}

	// check state is matched with state from cookie (against csrf attack)
	stateQuery := ctx.Query("state")
	log.Printf("get state: %s", stateQuery)

	if stateQuery != cookieState {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "state does not match"})
		return
	}

	// decode state
	stateJSON, err := base64.URLEncoding.DecodeString(stateQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode state"})
		return
	}

	var stateData map[string]string
	if err := json.Unmarshal(stateJSON, &stateData); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unmarshal state"})
		return
	}

	// check action
	action := stateData["action"]
	if action != "signup" && action != "login" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
		return
	}

	// get code
	code := ctx.Query("code")
	// googleConfig := configs.LoadAppConfig()

	// convert code to token
	token, err := googleConfig.Exchange(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert code to token"})
		return
	}

	// get user info from google api using token
	res, err := http.Get(conf.GoogleApiURL + token.AccessToken)
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

	// ユーザー情報をセッションに保存
	session := sessions.Default(ctx)
	session.Set("user_data", string(userData)) // userDataはsessionにinterface型として保存されるが、Getするした後string型で使用するのでstring型に変換している
	err = session.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user data in session"})
		return
	}
	log.Println("session:", session.Get("user_data"))

	// stateDataのactionでリダイレクト先を分岐
	switch action {
	case "signup":
		ctx.Redirect(http.StatusFound, conf.SignupRedirectURL)
	case "login":
		ctx.Redirect(http.StatusFound, conf.LoginRedirectURL)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
	}
}
