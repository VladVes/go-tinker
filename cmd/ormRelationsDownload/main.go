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
	err = db.AutoMigrate(&models.Comment{}, &models.Post{}, &models.User{})
	if err != nil {
		log.Fatalf("problem with automigrate profile model: %v ", err)
	}

	// -------------

	// ******************* Preload и Joins ***************************
	// Preload() — основной способ подгрузки связей.
	// Он делает дополнительные SELECT-запросы,
	// связывает данные по внешним ключам и аккуратно раскладывает их в структуры.
	// Joins() — более низкоуровневый инструмент: он строит SQL с JOIN
	// прямо на стороне базы и возвращает плоский результат.

	// +++++++++++++++++++++ Жадная загрузка с Preload: +++++++++++++++++
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

	// ------------ Joins() полезен, когда фильтровать нужно не связи,
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

	// =================== Отложенная и жадная загрузка =========================
	// По умолчанию GORM не подгружает связи «автоматом».
	// Это и есть отложенная (lazy) загрузка:
	// сначала в память попадает только основная модель,
	// а связи загружаются отдельно, когда это явно нужно.

	// ------ Отложенная загрузка через Association:
	var user models.User

	// Загружаем только пользователя, без постов.
	if err := db.First(&user, 1).Error; err != nil {
		log.Println("ошибка выборки пользователя:", err)
		return
	}
	// Posts пока пустой срез: данные не загружены.
	log.Println("постов до подгрузки:", len(user.Posts))

	// Отложенная загрузка постов через Association API.
	if err := db.
		Model(&user).
		Association("Posts").
		Find(&user.Posts); err != nil {
		log.Println("ошибка подгрузки постов:", err)
	}

	log.Println("постов после подгрузки:", len(user.Posts))
	// Сначала будет:
	// SELECT * FROM users
	// WHERE id = 1;
	// а при Association("Posts").Find — отдельный запрос:
	// SELECT * FROM posts
	// WHERE user_id = 1;

	// ------ Жадная (eager) загрузка делается через Preload():
	// связи приходят сразу, в рамках одного вызова GORM.:
	var user2 models.User

	// Жадная загрузка: пользователь и его посты одним вызовом ORM.
	if err := db.
		Preload("Posts").
		First(&user2, 1).Error; err != nil {
		log.Println("ошибка выборки:", err)
	}

	log.Println("постов:", len(user.Posts))

	// ------  Вложенные связи тоже можно подгружать жадно:
	// Жадная загрузка постов и комментариев к ним.
	if err := db.
		Preload("Posts.Comments").
		First(&user, 1).Error; err != nil {
		log.Println("ошибка выборки:", err)
	}

	// На практике удобно комбинировать подходы: список сущностей получать
	// с жадной подгрузкой (например, пользователей и их посты),
	// а тяжёлые связи (комментарии, историю изменений) подгружать отложенно
	// только там, где они действительно нужны.

	// =================== Ограничения и условия при загрузке связей ============
	//Preload() поддерживает условия и принимает вторым параметром
	// либо SQL-строку с аргументами, либо функцию,
	// которая настраивает вложенный *gorm.DB.

	var users4 []models.User

	// Подгружаем только опубликованные посты.
	if err := db.
		Preload("Posts", "published = ?", true).
		Find(&users4).Error; err != nil {
		log.Println("ошибка выборки:", err)
	}

	// Сортировка и ограничение количества подгружаемых записей через функцию:
	// Подгружаем не все посты, а только три последних для каждого пользователя.
	if err := db.
		Preload("Posts", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(3)
		}).
		Find(&users).Error; err != nil {
		log.Println("ошибка выборки:", err)
	}

	// Разным связям можно задать разные ограничения:

	// Подгружаем только опубликованные посты
	// и только одобренные комментарии.
	if err := db.
		Preload("Posts", "published = ?", true).
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("approved = ?", true)
		}).
		Find(&users).Error; err != nil {
		log.Println("ошибка выборки:", err)
	}

	// Условия работают и для вложенных связей:
	// Подгружаем только те теги, в имени которых есть «Go».
	if err := db.
		Preload("Posts.Tags", func(db *gorm.DB) *gorm.DB {
			return db.Where("tags.name LIKE ?", "%Go%")
		}).
		Find(&users).Error; err != nil {
		log.Println("ошибка выборки:", err)
	}

	// А если нужно отфильтровать основную сущность по связям
	// и в то же время аккуратно подгрузить связи — соединяем Joins() и Preload():
	// Берём только пользователей с опубликованными постами,
	// и сразу подгружаем эти посты.
	if err := db.
		Joins("JOIN posts ON posts.user_id = users.id AND posts.published = true").
		Preload("Posts", "published = ?", true).
		Find(&users).Error; err != nil {
		log.Println("ошибка выборки:", err)
	}
	// Joins() ограничит список пользователей,
	// Preload() сделает второй запрос и подгрузит к ним все опубликованные посты

	// Загрузка связанных данных в GORM — это управляемый процесс.
	// Preload() отвечает за жадную подгрузку связей,
	// Association — за отложенную загрузку по требованию,
	// Joins() — за фильтрацию основной сущности по связям.
	//
	// Условия и ограничения внутри Preload() позволяют не перетаскивать из базы лишнее.
	// Вся магия связей остаётся внутри ORM, а код работает с обычными структурами Go.

}
