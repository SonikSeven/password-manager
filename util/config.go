package util

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DbDriver         string `mapstructure:"DB_DRIVER"`
	DbSource         string `mapstructure:"DB_SOURCE"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDb       string `mapstructure:"POSTGRES_DB"`
	ServerAddress    string `mapstructure:"SERVER_ADDRESS"`
	JWTSecret        string `mapstructure:"JWT_SECRET"`
	EncryptionKey    string `mapstructure:"ENCRYPTION_KEY"`
	EncryptionKeyRaw []byte
}

func LoadConfig(path string) (cfg Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&cfg); err != nil {
		return
	}

	// decode the encryption key from Base64
	keyStr := strings.TrimSpace(cfg.EncryptionKey)
	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		log.Fatalf("failed to decode ENCRYPTION_KEY: %v", err)
	}
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		log.Fatalf("invalid AES key size: %d bytes", len(key))
	}

	cfg.EncryptionKeyRaw = key
	return
}
