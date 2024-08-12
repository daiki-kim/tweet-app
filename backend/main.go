package main

import (
	"github.com/gin-gonic/gin"

	"github.com/daiki-kim/tweet-app/backend/apps/controllers"
	"github.com/daiki-kim/tweet-app/configs"
)

func add(a, b int) int {
	return a + b
}

func main() {
	r := gin.Default()

	configs.InitializeAppConfig()

	r.POST("/google_login", controllers.GoogleLogin)
	r.POST("/google_callback", controllers.GoogleCallback)

	r.Run()
}
