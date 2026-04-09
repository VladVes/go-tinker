package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	case "update":
		handleUpdate(db, os.Args)
	case "show":
		handleShow(db, os.Args)
	case "delete":
		handleDelete(db, os.Args)
	default:
		log.Fatal("unknown actino")
	}
}

func handleList(db *gorm.DB) {
	movies := []models.Movie{}

	if err := db.Find(&movies).Error; err != nil {
		log.Fatalf("Error while getting movies list: %v", err)
	}

	for _, m := range movies {
		fmt.Println("--------------------------------------")
		fmt.Printf("Номер записи: %d\n", m.ID)
		fmt.Printf("Название:		%s\n", m.Title)
	}
}

const argsError = "[args] should be: <title> <gener> <released at: yyyy-mm-dd>"

func handleCreate(db *gorm.DB, args []string) {

	// Title        string `gorm:"size:150;not null;"`
	// Genre        string
	// ReleasedAt   time.Time
	// Descrtiption string
	// AdditionalInfo string
	// Rating float64 `gorm:"type:numeric(3,1)"`
	if len(args) < 6 {
		log.Fatal(argsError)
		return
	}

	title := args[3]
	gener := args[4]
	releasedAt, err := time.Parse("2006-01-02", args[5])
	if err != nil {
		log.Fatalf("parsint time error: %v", err)
	}
	log.Println(title, gener, releasedAt)

	movie := models.Movie{
		Title:      title,
		Genre:      gener,
		ReleasedAt: releasedAt,
	}

	if err := db.Create(&movie).Error; err != nil {
		log.Fatalf("movie creating error: %v", err)
	}
	log.Printf("Movie has been created successfuly! Movie record id: %d", movie.ID)

}

func handleUpdate(db *gorm.DB, args []string) {
	if len(args) < 6 {
		log.Fatal("[args] format: <id> <title|gener|releasedAt(yyyy-mm-dd)|description(text)|raiting(float))> <value>")
	}
	id, err := strconv.Atoi(args[3])
	if err != nil {
		log.Fatal("id should be numeric")
	}
	var movie models.Movie
	if err := db.First(&movie, id).Error; err != nil {
		log.Fatalf("error finding movie with id = %d: %v", id, err)
	}
	updateMovie := func(k string, v interface{}) {
		// интересно, что регистр значения в k не учитывается, т.е. к примерру
		// полуе в модели ReleasedAt а в метод Update передаётся "releasedAt" и всё работает
		if err := db.Model(&movie).Update(k, v).Error; err != nil {
			log.Fatalf("error updating %v", err)
		}
		// можно с исп. инструкции where, тогда не нужно предварительно находить сущность в БД с помо First:
		// if err := db.Model(&models.Movie{}).
		// 	Where("id = ?", args[3]).
		// 	Update(args[4], args[5]).Error; err != nil {
		// 	log.Fatal(err)
		// }
	}
	key := args[4]
	rawTextVal := args[5]
	switch key {
	case "releasedAt":
		{
			// В Go для форматирования даты и времени и их парсинга (преобразования строки в объект time.Time)
			// не используют стандартные обозначения типа yyyy-mm-dd —
			// вместо этого применяют «эталонную» (референтную) дату: Mon Jan 2 15:04:05 MST 2006.
			value, err := time.Parse("2006-01-02", rawTextVal)
			if err != nil {
				log.Fatalf("incorrect data format: %s. Error: %v", rawTextVal, err)
			}
			updateMovie(key, value)
		}
	case "rating":
		{
			value, err := strconv.ParseFloat(rawTextVal, 64)
			if err != nil {
				log.Fatalf("incorrect rating: %s. Error: %v", rawTextVal, err)
			}
			updateMovie(key, value)
		}
	default:
		updateMovie(key, rawTextVal)

	}
}

func handleShow(db *gorm.DB, args []string) {
	var movie models.Movie
	if len(args) < 4 {
		log.Fatal("error... need id")
	}
	id, err := strconv.Atoi(args[3])
	if err != nil {
		log.Fatalf("error pars id: %v", err)
	}
	if err := db.First(&movie, id).Error; err != nil {
		log.Fatalf("mvie with id = %d not found: %v", id, err)
	}
	fmt.Println(movie)
}

func handleDelete(db *gorm.DB, args []string) {
	if len(args) < 4 {
		log.Fatal("error... need id")
	}
	id, err := strconv.Atoi(args[3])
	if err != nil {
		log.Fatalf("error pars id: %v", err)
	}
	var movie models.Movie
	if err := db.Delete(&movie, id).Error; err != nil {
		log.Fatalf("movie with id = %d deletion err: %v", id, err)
	}

}
