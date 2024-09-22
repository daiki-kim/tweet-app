package controllers

import (
	"log"
	"net/http"

	"github.com/daiki-kim/tweet-app/backend/apps/dtos"
	"github.com/daiki-kim/tweet-app/backend/apps/services"
	utils "github.com/daiki-kim/tweet-app/backend/pkg"
	"github.com/gin-gonic/gin"
)

type ITweetController interface {
	CreateTweet(ctx *gin.Context)
	GetTweet(ctx *gin.Context)
	GetUserTweets(ctx *gin.Context)
}

type TweetController struct {
	service services.ITweetService
}

func NewTweetController(service services.ITweetService) ITweetController {
	return &TweetController{service: service}
}

func (c *TweetController) CreateTweet(ctx *gin.Context) {
	userId := getUserIdFromCtx(ctx)
	if userId == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert user id to uint64"})
		return
	}

	var input dtos.TweetInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
	}

	if err := c.service.CreateTweet(userId, input.Type, input.Content); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tweet"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *TweetController) GetTweet(ctx *gin.Context) {
	tweetId := getIdFromReq(ctx, "id")
	if tweetId == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert tweet id to uint64"})
		return
	}

	tweet, err := c.service.GetTweet(tweetId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tweet"})
		return
	}

	ctx.JSON(http.StatusOK, tweet)
}

func (c *TweetController) GetUserTweets(ctx *gin.Context) {
	userId := getIdFromReq(ctx, "user_id")
	if userId == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert user id to uint64"})
		return
	}

	tweets, err := c.service.GetUserTweets(userId)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user tweets"})
			return
		}
	}

	ctx.JSON(http.StatusOK, tweets)
}

// contextからstringのuser_idを取得してuintで返す
func getUserIdFromCtx(ctx *gin.Context) uint {
	userIdString, exist := ctx.Get("user_id")
	if !exist {
		return 0
	}
	userId := utils.String2Uint(userIdString.(string))
	return userId
}

// requestからstringのidを取得してuintで返す
func getIdFromReq(ctx *gin.Context, param string) uint {
	idString := ctx.Param(param)
	log.Println("[getIdFromReq] idString: ", idString)
	id := utils.String2Uint(idString)
	log.Println("[getIdFromReq] id: ", id)
	return id
}
