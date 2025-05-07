package postgres

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	const (
		Host     = "localhost"
		Port     = 5432
		User     = "postgres"
		Password = "postgres"
		DBName   = "wander_sphere"
		SSLMode  = "disable"
		TimeZone = "UTC"
	)

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		Host, Port, User, Password, DBName, SSLMode, TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	return db
}
