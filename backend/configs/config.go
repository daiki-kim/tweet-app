package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	GoogleLoginConfig oauth2.Config
}

var (
	AppConfig            Config
	GOOGLE_EMAIL_SCOPE   = "https://www.googleapis.com/auth/userinfo.email"
	GOOGLE_PROFILE_SCOPE = "https://www.googleapis.com/auth/userinfo.profile"
)

func GetEnvDefault(key, defVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defVal
	}

	return val
}

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func NewGoogleConfig() oauth2.Config {
	AppConfig.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  GetEnvDefault("GOOGLE_REDIRECT_URL", "http://localhost:8080/google/callback"),
		ClientID:     GetEnvDefault("GOOGLE_CLIENT_ID", ""),
		ClientSecret: GetEnvDefault("GOOGLE_CLIENT_SECRET", ""),
		Scopes:       []string{GOOGLE_EMAIL_SCOPE, GOOGLE_PROFILE_SCOPE},
		Endpoint:     google.Endpoint,
	}

	return AppConfig.GoogleLoginConfig
}

func InitializeAppConfig() {
	err := LoadEnv()
	if err != nil {
		log.Fatal("failed to load .env file: ", err)
	}

	NewGoogleConfig()
}
