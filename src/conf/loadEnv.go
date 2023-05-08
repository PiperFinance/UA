package conf

import (
	"github.com/spf13/viper"
	"time"
)

type config struct {
	DBHost         string `mapstructure:"DB_HOST"`
	DBUserName     string `mapstructure:"DB_USER"`
	DBUserPassword string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	DBPort         string `mapstructure:"DB_PORT"`

	MONGO_Url    string `mapstructure:"MONGO_URL"`
	MONGO_DBName string `mapstructure:"MONGO_DBNAME"`

	JwtAccessSecret     string        `mapstructure:"JWT_ACCESS_SECRET"`
	JwtRefreshSecret    string        `mapstructure:"JWT_REFRESH_SECRET"`
	JwtAccessExpiresIn  time.Duration `mapstructure:"JWT_ACCESS_EXPIRED_IN"`
	JwtRefreshExpiresIn time.Duration `mapstructure:"JWT_REFRESH_EXPIRED_IN"`
	JwtMaxAge           time.Duration `mapstructure:"JWT_MAXAGE"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`
}

var Config config

func LoadConfig(path string) error {
	conf := config{}

	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	//viper.SetConfigName("app")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&conf)
	Config = conf
	return err
}
