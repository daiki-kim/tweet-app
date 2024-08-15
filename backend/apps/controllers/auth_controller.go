package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/daiki-kim/tweet-app/backend/apps/dtos"
	"github.com/daiki-kim/tweet-app/backend/apps/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	SignupUsingOAuth(ctx *gin.Context)
	Signup(ctx *gin.Context)
}

type AuthController struct {
	service services.IAuthService
}

func NewAuthController(service services.IAuthService) IAuthController {
	return &AuthController{service: service}
}

// OAuthからのサインアップ
func (c *AuthController) SignupUsingOAuth(ctx *gin.Context) {
	// セッションからuserdataを取得
	session := sessions.Default(ctx)
	userData := session.Get("user_data")
	if userData == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to get user data from session"})
		return
	}

	// userdataをinputにバインド
	var input dtos.OAuthSignupInput
	if err := json.Unmarshal([]byte(userData.(string)), &input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unmarshal user data"})
		return
	}

	// 入力されたdobをctxからinputにバインド
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
	}

	// OAuthでGoogleからのユーザーデータを使用してサインアップ
	if err := c.service.SignupUsingOAuth(input.Name, input.Email, input.Dob); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to signup using OAuth"})
		return
	}

	// サインアップが成功したらセッションをクリア
	session.Delete("user_data")
	session.Save()

	ctx.Status(http.StatusCreated)
}

// Normalサインアップ
func (c *AuthController) Signup(ctx *gin.Context) {
	// NormalサインアップデータをDTOにバインド
	var input dtos.SignupInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
		return
	}

	// ユーザーデータからサインアップ
	if err := c.service.Signup(input.Name, input.Email, input.Dob, input.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to signup"})
		return
	}

	ctx.Status(http.StatusCreated)
}
