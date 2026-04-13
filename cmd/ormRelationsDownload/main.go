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

	if err := godotenv.Load("./cmd/ormRelationsDownload/.env"); err != nil {
		log.Printf("dotenv problem: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("no DATABASE_URL evn var, using default dsn value")
		dsn = defaultDsn
	}

	// --------------

	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // базовый вывод в консоль
		logger.Config{
			SlowThreshold: time.Second, // порог для медленных запросов
			LogLevel:      logger.Info, // подробный уровень логирования
			Colorful:      true,        // цветной вывод для удобства
		},
	)

	// --------------

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "app_",
		},
	})

	if err != nil {
		log.Fatalf("Gorm open problem: %v", err)
	}

	// -------------
	// err = db.AutoMigrate()
	// if err != nil {
	// 	log.Fatalf("problem with automigrate profile model: %v ", err)
	// }

	// -------------

	// ******************* Preload и Joins ***************************
	// Preload() — основной способ подгрузки связей.
	// Он делает дополнительные SELECT-запросы,
	// связывает данные по внешним ключам и аккуратно раскладывает их в структуры.
	// Joins() — более низкоуровневый инструмент: он строит SQL с JOIN
	// прямо на стороне базы и возвращает плоский результат.

	//---------- Жадная загрузка с Preload:
	var users []models.User

	// Загружаем всех пользователей и сразу подгружаем их посты.
	// пробуем в режиме dryrun:
	// tx := db.Session(&gorm.Session{
	// 	DryRun: true,
	// })

	// Формирование SELECT-запроса без обращения к базе
	// stmt := tx.Preload("Posts").Find(&users).Statement
	// log.Println(`Dry Run Preload("Posts").Find(&users): `, stmt.SQL.String())

	if err := db.Preload("Posts").Find(&users).Error; err != nil {
		log.Println("ошибка выборки:", err)
	}

	for _, u := range users {
		log.Println("пользователь:", u.Name, "постов:", len(u.Posts))
	}

	// ----------- Preload() умеет фильтровать связи:
	// Подгружаем только посты, в заголовке которых есть «Go».
	if err := db.
		Preload("Posts", "title LIKE ?", "%Go%").
		Find(&users).Error; err != nil {
		log.Println("ошибка выборки:", err)
	}

	// Joins() полезен, когда фильтровать нужно не связи,
	// а основную сущность по данным в связанных таблицах.
	var users2 []models.User

	// Выбираем только тех пользователей, у которых есть пост с GORM в заголовке.
	if err := db.
		Joins("JOIN posts ON posts.user_id = users.id").
		Where("posts.title LIKE ?", "%GORM%").
		Find(&users2).Error; err != nil {
		log.Println("ошибка выборки:", err)
	}

	// Частый приём: использовать Joins() для фильтрации,
	// а Preload() — для аккуратной подгрузки связей:
	// Здесь Joins() отберёт пользователей,
	// а Preload("Posts") подгрузит для них все посты (не только отфильтрованные).
	var users3 []models.User

	if err := db.
		Joins("JOIN posts ON posts.user_id = users.id").
		Where("posts.title LIKE ?", "%Go%").
		Preload("Posts").
		Find(&users3).Error; err != nil {
		log.Println("ошибка выборки:", err)
	}

}
