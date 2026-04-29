package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/VladVes/go-tinker/v2/internal/models"
)

const defaultDsn = "host=localhost user=postgres password=mysecretpassword dbname=postgres port=5432 sslmode=disable"

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
	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Product{})
	if err != nil {
		log.Fatalf("problem with automigrate profile model: %v ", err)
	}

	// Сырые запросы позволяют писать SQL напрямую,
	// но по-прежнему пользоваться параметрами и сканированием результатов в структуры.
	// gorm.Expr() встраивает отдельные выражения в обычные вызовы GORM
	// — для инкрементов, формул, CASE-условий, подзапросов и диалектных функций.

	// *************************** Raw и Exec *********************************
	// Raw() — для SELECT-ов (чтение и сканирование данных);
	// Exec() — для INSERT/UPDATE/DELETE и любых запросов без результата.

	// Оба метода принимают SQL-строку с ? и список параметров.
	// GORM передаёт их драйверу через database/sql, так что параметры не попадают
	// в SQL как конкатенация строк — защита от инъекций сохраняется.

	// -------------------- Raw: Чтение данных через Raw ----------------------
	var users models.User

	db.Raw("SELECT id, name FROM users WHERE name LIKE ?", "%Ан%").Scan(&users)
	// сформирует SQL:
	// SELECT
	// 		id,
	//		name
	// FROM users
	// WHERE name LIKE '%Ан%';

	// Точно так же можно читать агрегаты в простые типы:
	var count int64
	db.Raw("SELECT COUNT(*) FROM users WHERE active = ?", true).Scan(&count)
	log.Println("Активных пользователей:", count)

	// ------------------- Exec: Изменения через Exec ------------------------
	// Exec() используют, когда запрос ничего не возвращает
	// (или вас не интересует результат):

	result := db.Exec(
		"UPDATE users SET active = ? WHERE last_login < ?",
		false,
		time.Now().AddDate(0, -6, 0),
	)

	if result.Error != nil {
		log.Println("ошибка обновления:", result.Error)
	}

	log.Println("Обновлено строк:", result.RowsAffected)

	// Точно так же — массовое удаление:
	// result := db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	// log.Println("Удалено строк:", result.RowsAffected)

	// -------------- INSERT с RETURNING -------------------------------------
	// В PostgreSQL удобно сразу получить идентификатор свежесозданной записи
	// через RETURNING. Это тоже можно сделать через Raw():
	var id uint
	db.Raw(
		"INSERT INTO users (name) VALUES (?) RETURNING id",
		"Иван",
	).Scan(&id)

	log.Println("Создан пользователь с ID: ", id)
	// Здесь Raw() выполняет INSERT, а затем сканирует возвращённое значение в переменную.

	// --------------- Гибриды: сложные JOIN-ы и представления ----------------
	// Когда SELECT становится сложно выразить через Preload()/Joins(),
	// проще честно написать SQL, а результат положить в нужный тип:
	type PostWithAuthor struct {
		ID         uint
		Title      string
		Body       string
		AuthorName string
	}
	var postsWithAuthor []PostWithAuthor

	db.Raw(`
	SELECT p.id, p.title, p.body, u.name AS author_name
	FROM posts p
	JOIN users u ON p.user_id = u.id
	WHERE p.published = TRUE
	ORDER BY p.created_at DESC
	`).Scan(&postsWithAuthor)

	// ----------------- Exec() для процедур и функций -------------------------
	// Если база поддерживает хранимые процедуры или SQL-функции, Exec() подойдёт и для них:
	db.Exec("SELECT update_statistics(?)", "users")
	db.Exec("CALL recalc_totals(?)", 100)
	// Всё это происходит в том же контексте:
	// можно вызывать tx.Exec внутри транзакции,
	// добавлять WithContext() с тайм-аутами и видеть запросы в логгере GORM.

	// ******************* gorm.Expr и сложные ыражения *************************
	// Иногда нужен не целый кастомный запрос, а точечная «умная» операция:
	// инкремент счётчика без чтения, изменение поля по формуле, CASE-условие,
	// работа с JSONB, массивами, NOW() и другими функциями базы.
	// Для таких случаев у GORM есть gorm.Expr().

	// Expr() — это объект, который говорит:
	// «запиши здесь вот это SQL-выражение с параметрами».
	// Его используют в Update/Updates, Select(), Where(), Order()
	// и других местах, где обычно передаётся простое значение.

	// ------------------- Атомарный инкремент без чтения -----------------------
	db.Model(&models.Post{}).
		Where("id = ?", id).
		Update("views", gorm.Expr("COALESCE(views, 0) + ?", 1))
	// SQL:
	// UPDATE posts
	// SET views = COALESCE(views, 0) + 1
	// WHERE id = 42;

	// Инкремент выполняется на стороне базы, операция атомарная, гонок нет.

	// ----------------- Массовое обновление по формуле ------------------------
	// Допустим, нужно выдать 10% скидку на все книги.
	// Вместо того чтобы выгружать их в память и пересчитывать цену,
	// проще отдать формулу базе:
	db.Model(&models.Product{}).
		Where("category = ?", "Books").
		Update("price", gorm.Expr("ROUND(price * ?, 2)", 0.9))

	// ----------------- CASE в UPDATE: защита от отрицательных остатков--------
	// Ещё один популярный сценарий — не уходить в минус при списании со склада:
	n := 10
	db.Model(&models.Product{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"stock": gorm.Expr(
				"CASE WHEN stock >= ? THEN stock - ? ELSE stock END",
				n, n,
			),
			"updated_at": gorm.Expr("NOW()"),
		})
	// Логика «если достаточно товара — списать,
	// иначе оставить как есть» живёт в базе, а не в приложении.
	// Даже если два процесса обновят одну строку, CASE отработает корректно.

	// ------------ Вычисляемые поля в Select и Order -----------------------
	// Иногда удобнее сразу вернуть производное значение и сортировку по нему:
	type Row struct {
		Name      string
		AgeMonths int
	}

	var rows []Row

	db.Model(&models.User{}).
		Select("name, age * 12 AS age_months").
		Order(gorm.Expr("age * 12 DESC")).
		Scan(&rows)
	// База посчитает возраст в месяцах и отсортирует по нему же, GORM просто заберёт результат.

	// ---------- Диалектные фишки: JSONB, массивы, функции ------------------
	// GORM не знает всех операторов PostgreSQL, но Expr() позволяет аккуратно вставить их в запрос.

	// Например, обновление поля в JSONB:
	db.Model(&models.User{}).
		Where("id = ?", id).
		Update("data", gorm.Expr(
			"jsonb_set(COALESCE(data, '{}'::jsonb), '{profile,city}', to_jsonb(?), true)",
			"Москва",
		))
	// Или добавление тега в массив:
	db.Model(&models.Post{}).
		Where("id = ?", id).
		Update("tags", gorm.Expr("COALESCE(tags, '{}') || ARRAY[?]", "go"))

	// ********************* Подзапросы и upsert с выражениями *************************
	// Подзапрос в условии:
	type Customer struct {
		TotalSpent int
	}

	var customers []Customer

	avg := db.Table("orders").Select("AVG(total)")

	db.Model(&Customer{}).
		Where("total_spent > (?)", avg).
		Find(&customers)

	// Upsert с вычислениями при конфликте:
	db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "sku"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"stock":      gorm.Expr("products.stock + EXCLUDED.stock"),
			"updated_at": gorm.Expr("NOW()"),
		}),
	}).Create(&models.Product{SKU: "A1", Stock: 5})
	// PostgreSQL использует EXCLUDED.*, GORM аккуратно подставляет выражения в DO UPDATE.
	// несколько CTE (WITH / WITH RECURSIVE);
	// пачка оконных функций;
	// сложные HAVING, хинты оптимизатора, специфичные индексы;
	// многоуровневые JOIN-ы с ветвящейся логикой, удобнее честно перейти на «чистый»
	// SQL — через Raw() для SELECT и Exec() для изменений.
	// В больших проектах такие запросы часто выносят в отдельные файлы рядом
	// с миграциями или генерируют через инструменты вроде sqlc,
	// чтобы получить типобезопасные обёртки.
}
