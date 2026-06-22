package config

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	ConfigureR2CORS   bool
}

type CORSConfig struct {
	AllowedOrigins []string
}

func GetEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
func LoadDBConfig() *DBConfig {
	loadDotEnv()

	return &DBConfig{
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     GetEnv("DB_PORT", "3306"),
		User:     GetEnv("DB_USER", "exampleuser"),
		Password: GetEnv("DB_PASSWORD", "examplepassword"),
		Name:     GetEnv("DB_NAME", "coi"),
		SSLMode:  GetEnv("DB_SSLMODE", "disable"),
	}
}

func loadDotEnv() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func (c *DBConfig) MySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

func LoadS3Config() (*S3Config, error) {
	loadDotEnv()

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
		ConfigureR2CORS:   GetEnv("R2_CONFIGURE_BUCKET_CORS", "true") == "true",
	}
	return cfg, nil
}

func LoadCORSConfig() *CORSConfig {
	loadDotEnv()

	origins := GetEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000,http://localhost:5173,http://127.0.0.1:5173")
	return &CORSConfig{AllowedOrigins: splitCSV(origins)}
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}
