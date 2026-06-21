package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type S3Config struct {
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2BucketName      string
	R2Endpoint        string
	R2PublicBaseURL   string
}

func GetEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
func LoadDBConfig() *DBConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &DBConfig{
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     GetEnv("DB_PORT", "3306"),
		User:     GetEnv("DB_USER", "exampleuser"),
		Password: GetEnv("DB_PASSWORD", "examplepassword"),
		Name:     GetEnv("DB_NAME", "coi"),
		SSLMode:  GetEnv("DB_SSLMODE", "disable"),
	}
}

func (c *DBConfig) MySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

func LoadS3Config() (*S3Config, error) {
	accountID := GetEnv("ACCOUNT_ID", "")
	if accountID == "" {
		return nil, fmt.Errorf("Need R2Account id")
	}

	cfg := &S3Config{
		R2AccountID:       accountID,
		R2AccessKeyID:     GetEnv("ACCESS_KEY_ID", ""),
		R2SecretAccessKey: GetEnv("SECRET_ACCESS_KEY", ""),
		R2BucketName:      GetEnv("BUCKET_NAME", ""),
		R2Endpoint:        GetEnv("S3_ENDPOINT", ""),
		R2PublicBaseURL:   GetEnv("R2_PUBLIC_BASE_URL", ""),
	}
	return cfg, nil
}
