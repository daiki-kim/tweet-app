package dtos

type TweetInput struct {
	Type    string `json:"type" binding:"required,oneof=text image video"`
	Content string `json:"content" binding:"required,min=1,max=140"`
}
