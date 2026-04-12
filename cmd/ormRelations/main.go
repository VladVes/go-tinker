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

// GORM поддерживает четыре базовых типа связей, которые напрямую соответствуют
// привычным отношениям в базе данных:

// has one — «имеет одну» запись;
// has many — «имеет много»;
// belongs to — «принадлежит»;
// many2many — «многие ко многим».
// Связь задаётся комбинацией полей и тегов gorm.
// ORM анализирует структуру, находит внешние ключи,
// создаёт нужные столбцы и строит SQL при чтении и записи.

func main() {
	err := godotenv.Load("./cmd/ormRelations/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
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
	err = db.AutoMigrate(&models.Profile{}, &models.User{}, &models.Post{}, &models.Tag{})
	if err != nil {
		log.Fatalf("problem with automigrate profile model: %v ", err)
	}

	// ************ has one / belongs to - Связь один к одному ***************

	// Preload
	// Чтобы загрузить пользователя вместе с профилем, используется Preload():

	var user models.User

	// Загрузка пользователя и его профиля одним вызовом.
	if err := db.Preload("Profile").First(&user, 7).Error; err != nil {
		log.Println("ошибка выборки:", err)
	}

	log.Println("имя:", user.Name)
	log.Println("user obj", user)
	if user.Profile.ID != 0 {
		log.Println("bio:", user.Profile.Bio)
	}

	// ************ has many / belongs to - Связь один ко многим: has many / belongs to
	// Создаём пользователя.
	usr := models.User{Name: "Jane", Email: "Jane@mail.com"}
	if err := db.Create(&usr).Error; err != nil {
		log.Println("ошибка создания пользователя:", err)
	}

	// Привязываем посты через внешний ключ.
	posts := []models.Post{
		{Title: "Первый пост", Body: "Текст поста", UserID: usr.ID},
		{Title: "Второй пост", Body: "Ещё один текст", UserID: usr.ID},
	}
	if err := db.Create(&posts).Error; err != nil {
		log.Println("ошибка создания постов:", err)
	}

	// Для выборки автора вместе с его постами снова используется Preload():
	var u models.User
	result := db.Preload("Posts").First(&u, usr.ID)
	if result.Error != nil {
		log.Fatalf("preload user with posts problem: %v", err)
	}
	log.Println("автор:", u.Name, "количество постов:", len(u.Posts))

	for _, p := range u.Posts {
		log.Println("—", p.Title)
	}

	// ************ many2many - Связь многие ко многим ***********************
	// Связь многие ко многим: many2many

	// Связь многие ко многим возникает, когда сущности могутссылаться
	// друг на друга в обе стороны.
	// Например, у поста может быть несколько тегов, и
	// один и тот же тег встречается у разных постов.
	// В реляционной модели для этого заводят промежуточную таблицу —
	// таблицу связей

	// Создание поста с тегами выглядит так:

	post := models.Post{
		UserID: 8, // пользоватлье уже был создан с таким ID
		Title:  "Знакомство с GORM",
		Body:   "Пример работы с ORM в Go",
		Tags: []models.Tag{
			{Name: "Go"},
			{Name: "GORM"},
			{Name: "SQL"},
		},
	}

	if err := db.Create(&post).Error; err != nil {
		log.Println("ошибка создания поста:", err)
	}

	// Чтобы прочитать пост вместе с тегами, достаточно одного Preload():
	var p models.Post

	if err := db.Preload("Tags").First(&p, post.ID).Error; err != nil {
		log.Println("ошибка выборки поста:", err)
	}

	log.Println("пост:", p.Title)
	for _, tag := range p.Tags {
		log.Println("тег:", tag.Name)
	}

	// Association
	// Добавление нового тега к уже существующему посту делается через Associations() API:

	var p2 models.Post
	if err := db.First(&p2, post.ID).Error; err != nil {
		log.Println("пост не найден:", err)
		return
	}

	var tag models.Tag
	// Находим или создаём тег.
	if err := db.FirstOrCreate(&tag, models.Tag{Name: "Database"}).Error; err != nil {
		log.Println("ошибка с тегом:", err)
		return
	}

	// Привязываем тег к посту через таблицу post_tags.
	if err := db.Model(&p2).Association("Tags").Append(&tag); err != nil {
		log.Println("ошибка привязки тега:", err)
	}
	// GORM добавит новую строку в post_tags и обновит множество тегов у поста.

	// ************** Настройка внешних ключей и поведения при удалении **********************************
	err = db.AutoMigrate(&models.Order{})
	if err != nil {
		log.Fatalf("auto migration problem with Order model: %v", err)
	}
}
