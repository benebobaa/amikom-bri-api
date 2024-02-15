package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	PortApp                 string        `mapstructure:"PORT_APP"`
	DBDsn                   string        `mapstructure:"DB_DSN"`
	AppName                 string        `mapstructure:"APP_NAME"`
	GOEnv                   string        `mapstructure:"GO_ENV"`
	TokenAccessSymetricKey  string        `mapstructure:"TOKEN_ACCESS_SYMETRIC_KEY"`
	TokenRefreshSymetricKey string        `mapstructure:"TOKEN_REFRESH_SYMETRIC_KEY"`
	SecretKeyResetPassword  string        `mapstructure:"SECRET_KEY_RESET_PASSWORD"`
	TokenAccessDuration     time.Duration `mapstructure:"TOKEN_ACCESS_DURATION"`
	RefreshTokenDuration    time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	SecretKeyDuration       time.Duration `mapstructure:"SECRET_KEY_DURATION"`
	EmailName               string        `mapstructure:"EMAIL_NAME"`
	EmailSender             string        `mapstructure:"EMAIL_SENDER"`
	EmailPassword           string        `mapstructure:"EMAIL_PASSWORD"`
	OnesignalUrl            string        `mapstructure:"ONESIGNAL_URL"`
	OnesignalAppId          string        `mapstructure:"ONESIGNAL_APP_ID"`
	OnesignalAuthKey        string        `mapstructure:"ONESIGNAL_AUTH_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
