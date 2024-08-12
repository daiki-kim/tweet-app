package auth_test

import (
	"testing"
	"time"

	"github.com/daiki-kim/tweet-app/backend/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// テスト用の秘密鍵
var (
	testSignKey   = []byte("test_secret_key")
	testVerifyKey = []byte("test_secret_key")
)

// カスタムクレームからトークンが生成されるテスト
func TestGenerateTokenValidClaims(t *testing.T) {
	// テスト用のclaimを作成
	testCustomClaim := auth.NewClaim("test@example.com")

	// Act
	token, err := testCustomClaim.GenerateToken()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatalf("expected a valid token, got an empty string")
	}
}

// カスタムクレームからリフレッシュトークンが生成されるテスト
func TestGenerateRefreshTokenValidClaims(t *testing.T) {
	// テスト用のclaimを作成
	testCustomClaim := auth.NewClaim("test@example.com")

	// Act
	token, err := testCustomClaim.GenerateRefreshToken()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatalf("expected a valid token, got an empty string")
	}
}

// リフレッシュトークンをアップデートされるテスト
func TestUpdateRefreshToken(t *testing.T) {
	// テスト用のclaimを作成
	testCustomClaim := auth.NewClaim("test@example.com")

	// Act
	token, err := testCustomClaim.UpdateRefreshToken()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatalf("expected a valid token, got an empty string")
	}
}

// テスト用のカスタムクレームを持つトークンを生成する関数
func generateTestToken(email string, signKey []byte) (string, error) {
	claims := &auth.CustomClaim{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    auth.Issuer,
			Subject:   auth.Subject,
			Audience:  []string{auth.Audience},
			IssuedAt:  jwt.NewNumericDate(time.Now().Local()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(auth.TokenExpiration)),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signKey)
}

// 有効なトークンのテスト
func TestValidateToken_ValidToken(t *testing.T) {
	// テスト用のトークンを生成
	email := "test@example.com"
	tokenString, err := generateTestToken(email, testSignKey)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	// トークンを検証
	claims, err := auth.ValidateToken(tokenString, testVerifyKey)
	if err != nil {
		t.Fatalf("ValidateToken returned an error: %v", err)
	}

	// クレームが正しいか確認
	if claims.Email != email {
		t.Errorf("Expected email %v, got %v", email, claims.Email)
	}
}

// 無効な署名のトークンのテスト
func TestValidateToken_InvalidSignature(t *testing.T) {
	// テスト用のトークンを生成
	email := "test@example.com"
	tokenString, err := generateTestToken(email, []byte("wrong_secret_key"))
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	// トークンを検証
	_, err = auth.ValidateToken(tokenString, testVerifyKey)
	if err == nil {
		t.Fatal("Expected error due to invalid signature, got nil")
	}
}

// 期限切れトークンのテスト
func TestValidateToken_ExpiredToken(t *testing.T) {
	// 期限切れのトークンを生成
	claims := &auth.CustomClaim{
		Email: "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    auth.Issuer,
			Subject:   auth.Subject,
			Audience:  []string{auth.Audience},
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // 1時間前に期限切れ
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(testSignKey)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	// トークンを検証
	_, err = auth.ValidateToken(tokenString, testVerifyKey)
	if err == nil {
		t.Fatal("Expected error due to expired token, got nil")
	}
}

// 無効なトークンのテスト
func TestValidateToken_InvalidToken(t *testing.T) {
	// 不正なトークン文字列を使用
	invalidTokenString := "invalid.token.string"

	// トークンを検証
	_, err := auth.ValidateToken(invalidTokenString, testVerifyKey)
	if err == nil {
		t.Fatal("Expected error due to invalid token, got nil")
	}
}
