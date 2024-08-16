package main

import (
	"log"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/configs"
	"github.com/daiki-kim/tweet-app/backend/routes"
)

func main() {
	configs.InitializeConfig()
	err := models.SetDatabase(configs.Config.DBInstance)
	if err != nil {
		log.Fatal(err.Error())
	}

	db := models.DB
	r := routes.SetupRouter(db)

	r.Run()
}
