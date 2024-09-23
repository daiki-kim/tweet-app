package dtos

type FollowerInput struct {
	FolloweeID uint `json:"followee_id" binding:"required"`
}
