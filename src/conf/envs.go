package conf

import (
	"fmt"
	"net/url"
	"time"

	"github.com/caarlos0/env/v8"
	"go.uber.org/zap/zapcore"
)

type config struct {
	MongoUrl                     url.URL       `env:"MONGO_URL"`
	MongoDBName                  string        `env:"MONGO_DBNAME"`
	RedisUrl                     url.URL       `env:"REDIS_URL"`
	MongoSlowLoading             time.Duration `env:"MONGO_SLOW_LOAD" envDefault:"10s"`
	MongoMaxPoolSize             uint64        `env:"MONGO_MAX_POOL_SIZE" envDefault:"340"`
	MongoMaxTimeout              time.Duration `env:"MONGO_MAX_TIMEOUT" envDefault:"5m"`
	RedisDB                      int           `env:"REDIS_DB"`
	RedisMongoSlowLoading        time.Duration `env:"REDIS_SLOW_LOAD" envDefault:"3s"`
	DBHost                       string        `env:"DB_HOST"`
	DBUserName                   string        `env:"DB_USER"`
	DBUserPassword               string        `env:"DB_PASSWORD"`
	DBName                       string        `env:"DB_NAME"`
	DBPort                       string        `env:"DB_PORT"`
	MaxConcurrency               int           `env:"MAX_CONCURRENT_WORKER"`
	MaxTaskTimeOut               time.Duration `env:"MaxTaskTimeOut" envDefault:"10m"`
	AsynqMonUrl                  string        `env:"ASYNQ_MON_URL" envDefault:":7654"`
	ApiUrl                       string        `env:"API_URL" envDefault:":1300"`
	StartingBlockNumber          uint64        `env:"STARTING_BLOCK_NUMBER" envDefault:"3"`
	BlockHeadDelay               uint64        `env:"BLOCK_HEAD_DELAY" envDefault:"30"`
	SilenceRRCErrs               bool          `env:"RPC_ERROR_SILENCE" envDefault:"false"`
	SilenceParseErrs             bool          `env:"PARSE_ERROR_SILENCE" envDefault:"false"`
	SupportedChains              []int64       `env:"SUPPORTED_CHAINS" envSeparator:","`
	MultiCallTimeout             time.Duration `env:"PARSE_BLOCK_TIMEOUT" envDefault:"1m"`
	ParseBlockTimeout            time.Duration `env:"PARSE_BLOCK_TIMEOUT" envDefault:"2m"`
	FetchBlockTimeout            time.Duration `env:"FETCH_BLOCK_TIMEOUT" envDefault:"5m"`
	UserBalUpdateTimeout         time.Duration `env:"USER_BAL_UPDATE_TIMEOUT" envDefault:"5m"`
	TestTimeout                  time.Duration `env:"TEST_RPC_CONNECTION_TIMEOUT" envDefault:"15s"`
	ScanTaskTimeout              time.Duration `env:"SCAN_TASK_TIMEOUT" envDefault:"25s"`
	UpdateOnlineUsersTaskTimeout time.Duration `env:"ONLINE_USERS_TASK_TIMEOUT" envDefault:"25s"`
	LogLevel                     string        `env:"LOG_LEVEL" envDefault:"warn"`
	LogDir                       string        `env:"LOG_DIR" envDefault:"/var/bs/log"`
	MainnetDir                   string        `env:"MAINNET_DIR" envDefault:"/data/mainnets.json"`
	DEV                          bool          `env:"DEV_DEBUG" envDefault:"false"`
	LimitUsers                   bool          `env:"LIMIT_USERS" envDefault:"false"`
	JwtAccessSecret              string        `env:"JWT_ACCESS_SECRET"`
	JwtRefreshSecret             string        `env:"JWT_REFRESH_SECRET"`
	JwtAccessExpiresIn           time.Duration `env:"JWT_ACCESS_EXPIRED_IN"`
	JwtRefreshExpiresIn          time.Duration `env:"JWT_REFRESH_EXPIRED_IN"`
	JwtMaxAge                    time.Duration `env:"JWT_MAXAGE"`
	ClientOrigin                 string        `env:"CLIENT_ORIGIN"`
	UserAppURL                   url.URL       `env:"UA_URL" envDefault:"https://ua.piper.finance"`
	BlockScannerURL              url.URL       `env:"BS_URL" envDefault:"https://bs.piper.finance"`
	PortfolioScannerURL          url.URL       `env:"PS_URL" envDefault:"https://dev_ps.piper.finance"`
	PSSafeBalUrl                 string        `env:"PS_SAFE_BAL" envDefault:"/v1/tokens/balance/safe"`
	BSSetBalURL                  string        `env:"BS_SAFE_BAL" envDefault:"/bal/%d/"`
	PSSyncMinTimeDelay           time.Duration `env:"PS_SYNC_MIN_DELAY" envDefault:"10m"`

	ZapLogLevel zapcore.Level
}

var Config config

func LoadConfig() {
	if err := env.Parse(&Config); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	if lvl, err := zapcore.ParseLevel(Config.LogLevel); err != nil {
		panic(fmt.Sprintf("%+v", err))
	} else {
		Config.ZapLogLevel = lvl
	}
}
