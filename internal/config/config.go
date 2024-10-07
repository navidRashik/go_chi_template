package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel            string `mapstructure:"LOG_LEVEL"`
	SecretKey           string `mapstructure:"SECRET_KEY"`
	AllowedHosts        string `mapstructure:"ALLOWED_HOSTS"`
	ServerHost          string `mapstructure:"SERVER_HOST"`
	ServerPort          string `mapstructure:"SERVER_PORT"`
	MasterDbName        string `mapstructure:"MASTER_DB_NAME"`
	MasterDbUser        string `mapstructure:"MASTER_DB_USER"`
	MasterDbPassword    string `mapstructure:"MASTER_DB_PASSWORD"`
	MasterDbPort        string `mapstructure:"MASTER_DB_PORT"`
	MasterDbSslmode     string `mapstructure:"MASTER_SSL_MODE"`
	MasterDbHost        string `mapstructure:"MASTER_DB_HOST"`
	MasterDbUrl         string `mapstructure:"MASTER_DATABASE_URL"`
	SsoBase             string `mapstructure:"SSO_BASE_URL"`
	SsoClient           string `mapstructure:"SSO_CLIENT_ID"`
	SsoClientSecret     string `mapstructure:"SSO_CLIENT_SECRET"`
	Timeout             int    `mapstructure:"EXTERNAL_REQ_TIMEOUT"`
	TimeZone            string `mapstructure:"SERVER_TIMEZONE"`
	MockApiKey          string `mapstructure:"POSTMAN_MOCK_API_KEY"`
	FtBaseUrl           string `mapstructure:"FT_BASE_URL"`
	KafkaBrokerServers  string `mapstructure:"KAFKA_BROKER_SERVERS"`
	KafkaClientID       string `mapstructure:"KAFKA_CLIENT_ID"`
	KafkaGroupID        string `mapstructure:"KAFKA_GROUP_ID"`
	MftTopicName        string `mapstructure:"MFT_TOPIC_NAME"`
	DfcEmail            string `mapstructure:"DATA_FLUENT_LOGIN_EMAIL"`
	DfcPassword         string `mapstructure:"DATA_FLUENT_LOGIN_PASSWORD"`
	DfcBaseUrl          string `mapstructure:"DATA_FLUENT_LOGIN_BASE_URL"`
	NotificationBaseUrl string `mapstructure:"NOTIFICATION_BASE_URL"`
	RetryPeriod         int    `mapstructure:"RETRY_PERIOD"`
	RetryCount          int    `mapstructure:"RETRY_COUNT"`
}

var ConfigStruct Config

// SetupConfig configuration
func SetupConfig() (*Config, error) {
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&ConfigStruct); err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation(ConfigStruct.TimeZone)
	if err != nil {
		return nil, err
	}
	time.Local = loc
	// fmt.Printf("ConfigStruct: %v\n", ConfigStruct)
	return &ConfigStruct, nil
}

func (cfg *Config) GetDbDNSConfig() string {
	masterLocation := time.Local.String()
	masterDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.MasterDbHost, cfg.MasterDbUser, cfg.MasterDbPassword, cfg.MasterDbName, cfg.MasterDbPort, cfg.MasterDbSslmode, masterLocation,
	)
	return masterDSN
}

func (cfg *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
}

func init() {
	_, err := SetupConfig()
	if err != nil {
		panic(err)
	}
}
