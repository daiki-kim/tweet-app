package controllers

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/daiki-kim/tweet-app/backend/configs"
	utils "github.com/daiki-kim/tweet-app/backend/pkg"
)

func GoogleLogin(ctx *gin.Context) {
	configs.LoadAppConfig()
	conf := configs.Config

	state, err := utils.GenerateRandomString(32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}
	log.Printf("generate state: %s", state)
	log.Printf("ClientID: %s", conf.GoogleLoginConfig.ClientID)
	log.Printf("RedirectURL: %s", conf.GoogleLoginConfig.RedirectURL)

	ctx.SetCookie("oath_state", state, 3600, "/", "", false, true)

	// set redirect url with state
	url := conf.GoogleLoginConfig.AuthCodeURL(state)
	log.Println(url)
	ctx.Redirect(http.StatusFound, url)
}

func GoogleCallback(ctx *gin.Context) {
	configs.LoadAppConfig()
	conf := configs.Config

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
	log.Println(conf.GoogleApiURL)
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

	// サインアップページにリダイレクト
	ctx.Redirect(http.StatusFound, conf.SignupRedirectURL)
}
