package dtos

type TweetInput struct {
	Type    string `json:"type" binding:"required,oneof=text image video"`
	Content string `json:"content" binding:"required,min=1,max=140"`
}

type UpdateTweetInput struct {
	Type    string `json:"type" binding:"omitempty,oneof=text image video"`
	Content string `json:"content" binding:"omitempty,min=1,max=140"`
}
