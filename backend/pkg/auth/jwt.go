package auth

import (
	"errors"
	"time"

	"github.com/daiki-kim/tweet-app/backend/configs"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

// jwt constants
const (
	Subject                = "AccessToken"
	Issuer                 = "github.com/daiki-kim/tweet-app"
	Audience               = "github.com/daiki-kim/tweet-app"
	TokenExpiration        = time.Minute * time.Duration(10)
	RefreshTokenExpiration = time.Hour * time.Duration(1)
)

var (
	tokenSignKey          = []byte(configs.GetEnvDefault("TOKEN_SIGN_KEY", "secret"))
	TokenVerifyKey        = []byte(configs.GetEnvDefault("TOKEN_VERIFY_KEY", "secret"))
	refreshTokenSignKey   = []byte(configs.GetEnvDefault("REFRESH_TOKEN_SIGN_KEY", "secret"))
	refreshTokenVerifyKey = []byte(configs.GetEnvDefault("REFRESH_TOKEN_VERIFY_KEY", "secret"))
)

// custom claim struct with userId
type CustomClaim struct {
	UserId string
	jwt.RegisteredClaims
}

// create new claim with userId
func NewClaim(userId string) *CustomClaim {
	return &CustomClaim{
		UserId: userId,
	}
}

// generate new jwt token
func (c *CustomClaim) GenerateToken() (token string, err error) {
	// generate jwt standard token
	claims := jwt.RegisteredClaims{
		Issuer:    Issuer,
		Subject:   Subject,
		Audience:  []string{Audience},
		IssuedAt:  jwt.NewNumericDate(time.Now().Local()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(TokenExpiration)),
		ID:        uuid.New().String(),
	}

	// copy claims to custom claim
	if err := copier.CopyWithOption(c, claims, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return "", err
	}

	// generate jwt token
	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// add secret to jwt token
	token, err = generatedToken.SignedString(tokenSignKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generate new refresh token
func (c *CustomClaim) GenerateRefreshToken() (token string, err error) {
	// generate jwt standard token
	claims := jwt.RegisteredClaims{
		Issuer:    Issuer,
		Subject:   Subject,
		Audience:  []string{Audience},
		IssuedAt:  jwt.NewNumericDate(time.Now().Local()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(RefreshTokenExpiration)),
		ID:        uuid.New().String(),
	}

	// copy claims to custom claim
	if err := copier.CopyWithOption(c, claims, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return "", err
	}

	// generate jwt token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// add secret to jwt token
	token, err = refreshToken.SignedString(refreshTokenSignKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

// update refresh token
// TODO: if need to set expiration, use GenerateRefreshToken 2024-08-12
func (c *CustomClaim) UpdateRefreshToken() (token string, err error) {
	// generate jwt standard token
	claims := jwt.RegisteredClaims{
		Issuer:   Issuer,
		Subject:  Subject,
		Audience: []string{Audience},
	}

	// copy claims to custom claim
	if err := copier.CopyWithOption(c, claims, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return "", err
	}
	c.ID = uuid.New().String()
	c.IssuedAt = jwt.NewNumericDate(time.Now().Local())
	// c.ExpiresAt = jwt.NewNumericDate(time.Now().Local().Add(RefreshTokenExpiration))

	// generate updated jwt token
	updatedRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// add secret to updated jwt token
	token, err = updatedRefreshToken.SignedString(refreshTokenSignKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

// verify jwt token
// token: jwt token string (like: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c)
// parsedToken: jwt parsed token (like: *jwt.Token)
func ValidateToken(token string, verifyKey []byte) (*CustomClaim, error) {
	// parse jwt token
	// &CustomClaim{}: initialize custom claim and use it's pointer to parse token
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&CustomClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	// parse custom claim
	// *CustomClaim: use custom claim's pointer to parse token which are initialized in &CustomClaim{} in jwt.ParseWithClaims
	claims, ok := parsedToken.Claims.(*CustomClaim)
	if !ok || !parsedToken.Valid {
		err = errors.New("could not parse custom claims")
		return nil, err
	}

	return claims, nil
}
