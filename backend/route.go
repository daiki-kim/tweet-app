package main

import (
	"github.com/daiki-kim/tweet-app/backend/apps/controllers"
	"github.com/daiki-kim/tweet-app/backend/apps/repositories"
	"github.com/daiki-kim/tweet-app/backend/apps/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) *gin.Engine {
	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepository)
	authController := controllers.NewAuthController(authService)

	r := gin.Default()

	// セッションのミドルウェアを設定
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("my_session", store))

	apirouter := r.Group("/api")
	{
		v1Router := apirouter.Group("/v1")
		{
			authRouter := v1Router.Group("/signup")
			{
				authRouter.GET("/") // TODO: User data作成画面へredirect 2024-08-15
				authRouter.POST("/", authController.Signup)
				authRouter.GET("/oauth") // TODO: OAuthからのUser data作成画面へredirect 2024-08-15
				authRouter.POST("/oauth", authController.SignupUsingOAuth)
			}

		}
	}

	{
		r.GET("/google_login", controllers.GoogleLogin)
		r.GET("/google_callback", controllers.GoogleCallback)
	}

	return r
}
