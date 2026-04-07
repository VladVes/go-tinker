package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/VladVes/go-tinker/v2/internal/models"
)

const defaultDsn = "host=localhost user=postgres password=mysecretpassword dbname=postgres port=5432 sslmode=disable"

func main() {
	// CLI ожидает параметров
	if len(os.Args) < 3 {
		log.Fatal("usage: movies <list|create|show|update|delete> [args]")
	}

	err := godotenv.Load("./cmd/moviesCrudCLI/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = defaultDsn
	}

	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "app_",
		},
	})
	if err != nil {
		log.Fatalf("DB connection open error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("ошибка доступа к пулу соединений: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("ошибка пинга базы: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// -------------------- CLI -----------------------
	entity := os.Args[1]
	action := os.Args[2]
	if entity != "movies" {
		log.Fatal("only movies supported")
	}

	switch action {
	case "list":
		handleList(db)
	case "create":
		handleCreate(db, os.Args)
	default:
		log.Fatal("unknown actino")
	}
}

func handleList(db *gorm.DB) {
	movies := []models.Movie{}

	if err := db.Find(&movies).Error; err != nil {
		log.Fatalf("Error while getting movies list: %v", err)
	}

	for i, m := range movies {
		fmt.Println("--------------------------------------")
		fmt.Printf("Номер записи: %d\n", i)
		fmt.Printf("Название:		%s\n", m.Title)
		fmt.Printf("Жанр:   		%s\n", m.Genre)
		fmt.Printf("Дата выхода:	%s\n", m.ReleasedAt)
		fmt.Printf("Описание:		%s\n", m.Descrtiption)
		fmt.Printf("Рейтинг:		%.f\n", m.Rating)
	}
}

func handleCreate(db *gorm.DB, args []string) {
	movies := []models.Movie{}

	if err := db.Find(&movies).Error; err != nil {
		log.Fatalf("Error while getting movies list: %v", err)
	}

	// title := args[3]
	fmt.Println(args[3])
	fmt.Println(args[4])

	// TODO...

}
