package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type ConfigList struct {
	Env                 string
	DBInstance          int
	DBHost              string
	DBPort              int
	DBUser              string
	DBPassword          string
	DBName              string
	APICorsAllowOrigins []string

	GoogleLoginConfig oauth2.Config
	GoogleApiURL      string
	SignupRedirectURL string
	LoginRedirectURL  string
}

var (
	Config               ConfigList
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
	if err == nil {
		return nil
	}

	// テスト実行の時はテストがあるディレクトリの.envを読み込む
	err = godotenv.Load(".env.test")
	if err == nil {
		return nil
	}

	return err
}

func NewAppConfig() oauth2.Config {
	Config.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  GetEnvDefault("GOOGLE_REDIRECT_URL", "http://localhost:8080/google_callback"),
		ClientID:     GetEnvDefault("GOOGLE_CLIENT_ID", ""),
		ClientSecret: GetEnvDefault("GOOGLE_CLIENT_SECRET", ""),
		Scopes:       []string{GOOGLE_EMAIL_SCOPE, GOOGLE_PROFILE_SCOPE},
		Endpoint:     google.Endpoint,
	}

	return Config.GoogleLoginConfig
}

func LoadAppConfig() oauth2.Config {
	err := LoadEnv()
	if err != nil {
		log.Fatal("failed to load .env file: ", err)
	}

	return NewAppConfig()
}

func LoadConfig() error {
	err := LoadEnv()
	if err != nil {
		return err
	}

	DBInstance, err := strconv.Atoi(GetEnvDefault("DB_INSTANCE", "1"))
	if err != nil {
		return err
	}

	DBPort, err := strconv.Atoi(GetEnvDefault("DB_PORT", "3306"))
	if err != nil {
		return err
	}

	Config = ConfigList{
		Env:                 GetEnvDefault("ENV", "development"),
		DBInstance:          DBInstance,
		DBHost:              GetEnvDefault("DB_HOST", "0.0.0.0"),
		DBPort:              DBPort,
		DBUser:              GetEnvDefault("DB_USER", "app"),
		DBPassword:          GetEnvDefault("DB_PASSWORD", "password"),
		DBName:              GetEnvDefault("DB_NAME", "tweet_app"),
		APICorsAllowOrigins: []string{"http://0.0.0.0:8001"},

		GoogleLoginConfig: LoadAppConfig(),
		GoogleApiURL:      GetEnvDefault("GOOGLE_API_URL", "https://www.googleapis.com/oauth2/v3/userinfo"),
		SignupRedirectURL: GetEnvDefault("SIGNUP_REDIRECT_URL", "http://localhost:8080/api/v1/signup/oauth"),
		LoginRedirectURL:  GetEnvDefault("LOGIN_REDIRECT_URL", "http://localhost:8080/api/v1/login/oauth"),
	}

	return nil
}

func InitializeConfig() {
	if err := LoadConfig(); err != nil {
		log.Fatal("failed to load config: ", err)
	}
}
