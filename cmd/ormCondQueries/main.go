package main

import (
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
	if err := godotenv.Load("./cmd/ormCondQueries"); err != nil {
		log.Print("env file not found")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Print("env var DATABASE_URL not found, using default value")
		dsn = defaultDsn
	}

	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // базовый вывод в консоль
		logger.Config{
			SlowThreshold: time.Second, // порог для медленных запросов
			LogLevel:      logger.Info, // подробный уровень логирования
			Colorful:      true,        // цветной вывод для удобства
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
		log.Fatalf("DB connection: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Automigratino: %v", err)
	}
	log.Println("User automigration complete")

	// ++++++++++++++++++++++++++++++++ Where +++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// Where
	// Простейший пример: фильтр по возрасту.
	var users []models.User

	result := db.Where("age > ?", 20).Find(&users)
	if result.Error != nil {
		log.Println("ошибка выборки:", result.Error)
	}

	log.Println("найдено записей:", result.RowsAffected)

	// AND
	// Когда условий несколько
	age := 25
	city := "Москва"
	condQuery := db.Where("age = ?", age).Where("city = ?", city)
	result = condQuery.Find(&users)
	if result.Error != nil {
		log.Fatalf("find users error: %v", result.Error)
	}
	log.Println("Найдено записей: ", result.RowsAffected)
	log.Println(users)

}
