package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)
type DBConfig struct {
	Host	 string
	Port	 string
	User	 string
	Password string
	Name	 string
	SSLMode  string
}

func GetEnv(key, fallback string) string {
	if val:= os.Getenv(key); val!= ""{
		return val
	}
	return fallback
}
func LoadDBConfig() *DBConfig{
	if  err:= godotenv.Load();err!=nil{
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &DBConfig{
		Host: GetEnv("DB_HOST","localhost"),
		Port: GetEnv("DB_PORT","3306"),
		User: GetEnv("DB_USER","exampleuser"),
		Password: GetEnv("DB_PASSWORD","examplepassword"),
		Name: GetEnv("DB_NAME","coi"),
		SSLMode: GetEnv("DB_SSLMODE","disable"),
	}
}

func (c *DBConfig) MySQLDSN() string{
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	c.User, c.Password, c.Host, c.Port, c.Name)
}