package routes

import (
	"github.com/daiki-kim/tweet-app/backend/apps/controllers"
	"github.com/daiki-kim/tweet-app/backend/apps/repositories"
	"github.com/daiki-kim/tweet-app/backend/apps/services"
	"github.com/daiki-kim/tweet-app/backend/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepository)
	authController := controllers.NewAuthController(authService)

	tweetRepository := repositories.NewTweetRepository(db)
	tweetService := services.NewTweetService(tweetRepository)
	tweetController := controllers.NewTweetController(tweetService)

	followerRepository := repositories.NewFollowerRepository(db)
	followerService := services.NewFollowerService(followerRepository)
	followerController := controllers.NewFollowerController(followerService)

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
				signupRouter.GET("/")                                        // TODO: Frontend実装後にUser data作成画面へredirectする 2024-08-15
				signupRouter.POST("/", authController.Signup)                // User data送信先
				signupRouter.GET("/oauth")                                   // TODO: Frontend実装後にOAuthからのUser data作成画面へredirectする 2024-08-15
				signupRouter.POST("/oauth", authController.SignupUsingOAuth) // OAuth経由のUser data送信先
			}

			loginRouter := v1Router.Group("/login")
			{
				loginRouter.GET("/")                                      // TODO: User data入力画面へredirect 2024-08-15
				loginRouter.POST("/", authController.Login)               // User data送信先
				loginRouter.GET("/oauth", authController.LoginUsingOAuth) // OAuthからのリダイレクト先
			}

			tweetRouterWithAuth := v1Router.Group("/tweet", middlewares.JwtTokenVerifier())
			{
				tweetRouterWithAuth.POST("/", tweetController.CreateTweet)
				tweetRouterWithAuth.GET("/:id", tweetController.GetTweet)
				tweetRouterWithAuth.GET("/user/:user_id", tweetController.GetUserTweets)
				tweetRouterWithAuth.PUT("/:id", tweetController.UpdateTweet)
				tweetRouterWithAuth.DELETE("/:id", tweetController.DeleteTweet)
			}

			followerRouterWithAuth := v1Router.Group("/follower", middlewares.JwtTokenVerifier())
			{
				followerRouterWithAuth.POST("/", followerController.Follow)
				followerRouterWithAuth.GET("/:id", followerController.GetFollower)
				followerRouterWithAuth.GET("/followers/:followee_id", followerController.GetFollowers) // followee_idのユーザーをフォローしているユーザ一覧を取得
				followerRouterWithAuth.DELETE("/:id", followerController.DeleteFollower)
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
