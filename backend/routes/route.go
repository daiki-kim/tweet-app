package routes

import (
	"github.com/daiki-kim/tweet-app/backend/apps/controllers"
	"github.com/daiki-kim/tweet-app/backend/apps/repositories"
	"github.com/daiki-kim/tweet-app/backend/apps/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
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
			signupRouter := v1Router.Group("/signup")
			{
				signupRouter.GET("/")                                        // TODO: User data作成画面へredirect 2024-08-15
				signupRouter.POST("/", authController.Signup)                // User data送信先
				signupRouter.GET("/oauth")                                   // TODO: OAuthからのUser data作成画面へredirect 2024-08-15
				signupRouter.POST("/oauth", authController.SignupUsingOAuth) // OAuth経由のUser data送信先
			}

			loginRouter := v1Router.Group("/login")
			{
				loginRouter.GET("/")                                      // TODO: User data入力画面へredirect 2024-08-15
				loginRouter.POST("/", authController.Login)               // User data送信先
				loginRouter.GET("/oauth", authController.LoginUsingOAuth) // OAuthからのリダイレクト先
			}

		}
	}

	{
		r.GET("/google_login/:action", func(ctx *gin.Context) {
			action := ctx.Param("action")
			controllers.GoogleLogin(ctx, action)
		})
		r.GET("/google_callback", controllers.GoogleCallback)
	}

	return r
}
