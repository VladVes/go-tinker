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
	if err := godotenv.Load("./cmd/lessons/orm-cond-queries/.env"); err != nil {
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

	// AND ( Where().Where().Where() )
	// Когда условий несколько
	age := 25
	city := "Москва"
	// SQL: WHERE age > 25 AND city = 'Москва'
	condQuery := db.Where("age = ?", age).Where("city = ?", city)
	result = condQuery.Find(&users)
	if result.Error != nil {
		log.Fatalf("find users error: %v", result.Error)
	}
	log.Println("Найдено записей: ", result.RowsAffected)
	log.Println(users)

	// Or ( Where().Or() )
	// Or() добавляет условие, которое расширяет выборку.
	// SQL: WHERE city = 'Москва' OR age < 18
	query := db.Where("city = ?", "Москва").Or("age < ?", 18)
	result = query.Find(&users)
	if result.Error != nil {
		log.Println("ошибка выборки:", result.Error)
	}

	log.Println(`
	query := db.Where("city = ?", "Москва").Or("age < ?", 18)
	Найдено записей:`, result.RowsAffected)
	log.Println(users)

	// Not - Исключение значений
	// SQL: WHERE NOT (city = 'Москва')
	// Все, кроме москвичей.
	result = db.Not("city = ?", "Москва").Find(&users)
	if result.Error != nil {
		log.Println("ошибка выборки:", result.Error)
	}
	log.Println(`
	result = db.Not("city = ?", "Москва").Find(&users)
	Найдено записей:`, result.RowsAffected)
	log.Println(users)
	// NOT c несколькими значениями передаваемыми в map
	conditions := map[string]any{
		"city": []string{"Москва", "Питер"},
	}
	result = db.Not(conditions).Find(&users)
	if result.Error != nil {
		log.Println("ошибка выборки:", result.Error)
	}

	log.Println(`
	conditions := map[string]any{
		"city": []string{"Москва", "Питер"},
	}

	result = db.Not(conditions).Find(&users)

	Найдено записей:`, result.RowsAffected)
	log.Println(users)

	//Where() принимает не только строки с ?, но и карты:
	// SQL: WHERE city = 'Москва' AND age = 30
	whereConditions := map[string]any{
		"city": "Москва",
		"age":  30,
	}
	// Если значение — срез, ORM строит IN.
	// Структуры работают аналогично, но используют имена полей Go и теги gorm:"column:...".

	// Подводный камень: нулевые значения в структурах
	// Если фильтрация должна учитывать 0, пустую строку или false, структура не подходит:
	// GORM пропускает такие поля.
	// Именно поэтому структуры редко используют для динамических фильтров — проще карта.

	if err := db.Where(whereConditions).Find(&users).Error; err != nil {
		log.Println("ошибка выборки:", err)
	}

	// ++++++++++++++++++++++++++++++++ Select +++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// По умолчанию GORM выбирает все столбцы. Но часто нужны только несколько полей:
	// SQL:
	// 	SELECT
	//     name,
	//     email
	// FROM users
	// WHERE age > 20;
	// Это ускоряет запросы и уменьшает объём данных.
	query = db.Select("login", "email").Where("age > ?", 20)
	result = query.Find(&users)

	if result.Error != nil {
		log.Println("ошибка выборки:", result.Error)
	}
	log.Println(`
	query = db.Select("name", "email").Where("age > ?", 20)
	result = query.Find(&users)

	Найдено записей:`, result.RowsAffected)
	log.Println(users)
	// Опасный момент: вычисляемые поля в модель User загружать нельзя
	// Иногда хочется посчитать что-то «на лету»: query := db.Select("name, age * 12 AS age_in_months").Where("age > ?", 18)
	// Такой запрос нельзя мапить в модель User

	// ++++++++++++++++++++++++ Частые выражения: IN, LIKE, BETWEEN ++++++++++++++++++++++++++++++++++
	// IN:
	cities := []string{"Москва", "Питер", "Казань"}
	db.Where("city IN ?", cities).Find(&users)

	// LIKE - Поиск подстроки:
	pattern := "%ан%"
	db.Where("LOWER(name) LIKE LOWER(?)", pattern).Find(&users)

	// BETWEEN - Диапазоны:
	db.Where("age BETWEEN ? AND ?", 18, 30).Find(&users)
	// Все параметры передаются безопасно — никакой конкатенации строк.

	// так же можно использовать различные имеющиеся стандартные механизмы для сортировки и ограничения выдачи
	// к примеру:
	// Order("created_at DESC").Limit(10)

	// Плохой вариант:

	// "age > " + strconv.Itoa(age)
	// Такие строки:

	// плохо читаются,
	// легко ломаются,
	// небезопасны (SQL-инъекции).
	// Правильный способ — WHERE age > ? или карта условий.

	// 	Все вызовы Where(), Or(), Not(), Select() попадают в объект Statement. Перед выполнением GORM:
	// собирает условия в единый SQL-текст,
	// подставляет ?,
	// передаёт параметры отдельно в driver.
}
