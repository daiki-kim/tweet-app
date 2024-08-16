package dtos

type OAuthSignupInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Dob   string `json:"dob" binding:"required"`
}

type SignupInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"min=8"`
	Dob      string `json:"dob" binding:"required"`
}

type OAuthLoginInput struct {
	Email string `json:"email" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
