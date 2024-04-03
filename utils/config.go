package utils

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	AWS_BUCKET_NAME                 string `mapstructure:"AWS_BUCKET_NAME"`
	AWS_REGION                      string `mapstructure:"AWS_REGION"`
	AWS_S3_BUCKET_ACCESS_KEY        string `mapstructure:"AWS_S3_BUCKET_ACCESS_KEY"`
	AWS_S3_BUCKET_SECRET_ACCESS_KEY string `mapstructure:"AWS_S3_BUCKET_SECRET_ACCESS_KEY"`
}

func LoadViperEnvironment(path string) (config Config, err error) {
	viper.SetConfigFile(filepath.Join(path, ".env"))
	viper.SetConfigType("env")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return config, nil
}
