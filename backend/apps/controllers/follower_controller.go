package controllers

import (
	"net/http"

	"github.com/daiki-kim/tweet-app/backend/apps/dtos"
	"github.com/daiki-kim/tweet-app/backend/apps/services"
	"github.com/gin-gonic/gin"
)

type IFollowerController interface {
	Follow(ctx *gin.Context)
	GetFollower(ctx *gin.Context)
	GetFollowers(ctx *gin.Context)
	DeleteFollower(ctx *gin.Context)
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

func (c *FollowerController) DeleteFollower(ctx *gin.Context) {
	userId := getUserIdFromCtx(ctx)
	if userId == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user id"})
		return
	}

	id := getIdFromReq(ctx, "id")
	if id == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get follower id"})
		return
	}

	if err := c.service.DeleteFollower(id, userId); err != nil {
		if err.Error() == "follower not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "you don't have permission to delete this follower" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete follower"})
		}
	}

	ctx.Status(http.StatusOK)
}

func (c *FollowerController) GetFollower(ctx *gin.Context) {
	id := getIdFromReq(ctx, "id")
	if id == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get follower id"})
		return
	}

	follower, err := c.service.GetFollower(id)
	if err != nil {
		if err.Error() == "follower not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get follower"})
			return
		}
	}

	ctx.JSON(http.StatusOK, follower)
}
