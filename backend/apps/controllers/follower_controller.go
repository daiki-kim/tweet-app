package controllers

import (
	"net/http"

	"github.com/daiki-kim/tweet-app/backend/apps/dtos"
	"github.com/daiki-kim/tweet-app/backend/apps/services"
	"github.com/gin-gonic/gin"
)

type IFollowerController interface {
	Follow(ctx *gin.Context)
	GetFollowers(ctx *gin.Context)
}

type FollowerController struct {
	service services.IFollowerService
}

func NewFollowerController(service services.IFollowerService) IFollowerController {
	return &FollowerController{service: service}
}

func (c *FollowerController) Follow(ctx *gin.Context) {
	followerId := getUserIdFromCtx(ctx)
	if followerId == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get follower id"})
		return
	}

	var followerInput dtos.FollowerInput
	if err := ctx.ShouldBindJSON(&followerInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
		return
	}

	follower, err := c.service.Follow(followerId, followerInput.FolloweeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to follow"})
		return
	}

	ctx.JSON(http.StatusCreated, follower)
}

func (c *FollowerController) GetFollowers(ctx *gin.Context) {
	followeeId := getIdFromReq(ctx, "followee_id")
	if followeeId == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get followee id"})
		return
	}

	followersUserData, err := c.service.GetFollowers(followeeId)
	if err != nil {
		if err.Error() == "followers not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get followers"})
			return
		}
	}

	ctx.JSON(http.StatusOK, followersUserData)
}
