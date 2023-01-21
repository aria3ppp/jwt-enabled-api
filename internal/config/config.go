package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var Config *config = initConfig("config.yml")

type config struct {
	Postgres struct {
		DB       string `yaml:"db" env:"POSTGRES_DB" env-required:"true"`
		User     string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
		Host     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
		Port     uint16 `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	} `yaml:"postgres" env-required:"true"`

	Auth struct {
		ECDSASigningKeyBase64 string `yaml:"ecdsa_signing_key_base64" env:"ECDSA_SIGNING_KEY_BASE64" env-required:"true"`
		ExpireInSecs          struct {
			Jwt     int `yaml:"jwt" env-required:"true"`
			Refresh int `yaml:"refresh" env-required:"true"`
		} `yaml:"expire_in_secs" env-required:"true"`
	} `yaml:"auth" env-required:"true"`
}

func initConfig(filename string) *config {
	// find config file path
	_, currentFile, _, _ := runtime.Caller(0)
	pathDir := filepath.Join(filepath.Dir(currentFile), "..", "..")
	// load .env file
	if err := godotenv.Load(filepath.Join(pathDir, ".env")); err != nil {
		log.Fatalf("failed loading .env file: %s", err)
	}
	// read configs
	var config config
	if err := cleanenv.ReadConfig(filepath.Join(pathDir, filename), &config); err != nil {
		log.Fatalf("config.initConfig: failed loading configs: %s", err)
	}
	return &config
}
