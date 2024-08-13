package main

import (
	"github.com/gin-gonic/gin"

	"github.com/daiki-kim/tweet-app/backend/apps/controllers"
	"github.com/daiki-kim/tweet-app/backend/configs"
)

func main() {
	configs.InitializeConfig()

	r := gin.Default()

	configs.LoadAppConfig()

	r.GET("/google_login", controllers.GoogleLogin)
	r.GET("/google_callback", controllers.GoogleCallback)

	r.Run()
}
