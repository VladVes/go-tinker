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

	err := godotenv.Load("./cmd/movies-crud-cli/.env")
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

	// --------------------------------------------------------------------------------
	// --------------------- Relatinos example -----------------------------------

	err = db.AutoMigrate(&models.Review{}, &models.Movie{}, &models.Director{}, &models.Actor{})
	if err != nil {
		log.Fatalf("auto migration problem: %v", err)
	}
	log.Println("auto migragtion complete")

	newDirector := models.Director{
		Name: "Louis Armstrong",
	}

	if err := db.FirstOrCreate(&newDirector).Error; err != nil {
		log.Fatalf("new director creation: %v", err)
	}

	releasedAt, _ := time.Parse("2006-01-02", "2026-12-04")
	newMovie := models.Movie{
		Title:       "Super movie",
		Genre:       "sci-fi",
		ReleasedAt:  releasedAt,
		Description: "new film",
		DirectorID:  newDirector.ID,
		Actor: []models.Actor{
			{Name: "John Doe"},
			{Name: "Alеizee Jorry"},
			{Name: "Nikolay Bulkin"},
		},
	}

	// var m models.Movie
	// FirstOrCreate не создал записей в Actor
	// if err := db.FirstOrCreate(&m, &newMovie).Error; err != nil {
	// 	log.Fatalf("test movie creation: %v", err)
	// }
	movieResult := db.Create(&newMovie)
	if movieResult.Error != nil {
		log.Fatalf("Problem with test movie creation: %v", movieResult.Error)
	}
	log.Printf("new test movie with director and actros has been created!: %v", newMovie.ID)

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
	case "unraited":
		handleUnrated(db)
	case "review":
		handleAddReview(db, os.Args)
	case "rating":
		handleGetRating(db, os.Args)
	default:
		log.Fatal("unknown actino")
	}
}

func handleList(db *gorm.DB) {
	movies := []models.Movie{}

	if err := db.Preload("Director").Find(&movies).Error; err != nil {
		log.Fatalf("Error while getting movies list: %v", err)
	}

	for _, m := range movies {
		fmt.Println("--------------------------------------")
		fmt.Printf("Номер записи: %d\n", m.ID)
		fmt.Printf("Название:		%s\n", m.Title)
		fmt.Printf("Режиссёр:		%s\n", m.Director.Name)
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
	updateMovie := func(k string, v any) {
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
	if err := db.Preload("Actor").Preload("Director").First(&movie, id).Error; err != nil {
		log.Fatalf("mvie with id = %d not found: %v", id, err)
	}
	fmt.Println("Результат: ")
	fmt.Printf("Номер записи %d\n", movie.ID)
	fmt.Printf("Название: %s\n", movie.Title)
	fmt.Printf("Жанр: %s\n", movie.Genre)
	fmt.Printf("Описание: %s\n", movie.Description)
	fmt.Printf("Режиссёр: %v\n", movie.Director.Name)
	fmt.Println("Актёры:")
	for _, a := range movie.Actor {
		fmt.Println(a.Name)
	}
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

func handleUnrated(db *gorm.DB) {
	var movies []models.Movie
	fieldsToSelect := []string{"id", "title"}
	result := db.Select(fieldsToSelect).Where("rating IS NULL").Find(&movies)
	if result.Error != nil {
		log.Fatalf("getting movies with nil rating: %v", result.Error)
	}

	for _, m := range movies {
		log.Printf("№ %d\n", m.ID)
		log.Println(m.Title)
	}

}

// ----- Проверки создания отзыва ---------
// Пример проверок на уровне приложения
// func validateReview(score int, text string) error {
// 	if text == "" {
// 		return errors.New("review text is required")
// 	}
// 	if score < 1 || score > 10 {
// 		return errors.New("score must be between 1 and 10")
// 	}
// 	return nil
// }

// Пример проверки на уровне базы
// CREATE UNIQUE INDEX IF NOT EXISTS reviews_movie_text_idx
// ON reviews (movie_id, text);
// ----- Пример как создаются ошибки для их использованию
// var (
// 	errEmptyTitle       = errors.New("title is empty")
// 	errTooLongTitle     = errors.New("title is too long")
// 	errMissingAssignee  = errors.New("assignee is required for done tasks")
// 	errCreateDone       = errors.New("cannot create task with done status")
// 	errProjectNotFound  = errors.New("project not found")
// 	errAssigneeNotFound = errors.New("assignee not found")
// )

func handleAddReview(db *gorm.DB, args []string) {
	if len(args) < 6 {
		log.Fatalf("Недостаточно данных")
	}

	movieID, err := strconv.ParseUint(args[3], 10, 64)
	if err != nil {
		log.Fatal("parsing movieID err: ", err)
	}
	score, err := strconv.Atoi(args[4])
	if err != nil {
		log.Fatal("parsing score err: ", err)
	}

	review := models.Review{
		MovieID: uint(movieID),
		Score:   score,
		Text:    args[5],
	}

	// var movie models.Movie

	db.Transaction(func(tx *gorm.DB) error {
		// Пример обработки ошибки дубликата
		// if err := db.Create(&review).Error; err != nil {
		// 	if strings.Contains(err.Error(), "reviews_movie_text_idx") {
		// 		return errors.New("duplicate review for this movie")
		// 	}
		// 	return err
		// }
		if err := tx.Create(&review).Error; err != nil {
			log.Fatalf("create review problem: %v", err)
			return err
		}

		// if err := tx.First(&movie, movieID).Error; err != nil {
		// 	log.Fatalf("create review problem - no movie found: %v", err)
		// 	return err
		// }

		if err := tx.Model(&models.Movie{}).
			Where("id = ?", movieID).
			Update("reviewCount", gorm.Expr("review_count + ?", 1)).Error; err != nil {
			log.Fatalf("update movie review counter problem: %v", err)
			return err
		}
		return nil
	})
}

func handleGetRating(conn *gorm.DB, args []string) {
	if len(args) < 3 {
		log.Println(args)
		log.Fatalf("Недостаточно данных")
	}

	type MovieRatingDTO struct {
		Title        string
		Rating       float64
		ReviewsCount int64
	}
	var movieRatingDTO []MovieRatingDTO

	conn.Raw(`
SELECT
    m.title,
    AVG(r.score) AS rating,
    COUNT(r.id) AS reviews_count
FROM app_movies m
LEFT JOIN app_reviews r ON r.movie_id = m.id
GROUP BY m.id, m.title
ORDER BY rating DESC NULLS LAST, m.title;
	`).Scan(&movieRatingDTO)

	log.Println(movieRatingDTO)
}
