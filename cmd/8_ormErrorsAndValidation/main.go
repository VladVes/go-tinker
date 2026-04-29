package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
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

}

// ***********************************************************************************
// как интерпретировать ошибки GORM (ErrRecordNotFound, ErrInvalidData и не только);
// как проверять входные данные до запроса к базе;
// как аккуратно работать с ошибками миграций и транзакций, чтобы база всегда оставалась в согласованном состоянии.

// =================== Ошибки ErrRecordNotFound и ErrInvalidData =====================
// GORM почти всегда возвращает результат операции через объект *gorm.DB. В нём есть:
// Error — сама ошибка;
// RowsAffected — сколько строк реально затронуто.

// gorm.ErrRecordNotFound - специальная ошибка, когда выборка ничего не нашла
// gorm.ErrInvalidData - когда передали некорректную модель

// ----------- ErrRecordNotFound: Отличаем «не найдено» от «сломалось» -----
// Пример аккуратного чтения записи:
func DemoErrRecordNotFound(db *gorm.DB) (*models.User, error) {
	var user models.User
	email := "alice@mail.com"

	res := db.Where("email = ?", email).First(&user)

	switch {
	case errors.Is(res.Error, gorm.ErrRecordNotFound):
		// доменная «не найдено» — можно вернуть 404 или пустой ответ
		return nil, fmt.Errorf("user not found: %w", res.Error)
	case res.Error != nil:
		// что-то пошло не так на уровне БД или драйвера
		return nil, fmt.Errorf("db query failed: %w", res.Error)
	default:
		// запись найдена
		return &user, nil
	}
}

// для сравнения используем errors.Is(), а не сравнение текстов сообщений;
// ошибку оборачиваем через %w, чтобы верхний уровень всё ещё мог узнать,
// что внутри был ErrRecordNotFound.

// ------------ ErrRecordNotFound и RowsAffected при массовых операциях -----------------------
// Если вы делаете массовый UPDATE или DELETE, ErrRecordNotFound не появится — запрос корректен,
// просто не нашёл подходящих строк. В этом случае смотрят на RowsAffected:

func DemoRowsAffected(db *gorm.DB) error {
	res := db.Model(&models.Order{}).
		Where("status = ?", "draft").
		Update("status", "submitted")

	if res.Error != nil {
		return fmt.Errorf("update orders failed: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		// по смыслу — «не найдено ни одной строки под условие»
		return fmt.Errorf("no orders updated: %w", gorm.ErrRecordNotFound)
	}

	return nil
}

// Такой код явно фиксирует ситуацию «операция ничего не изменила» и
// может превратить её в 422/404, а не тихо делать вид, что всё прошло успешно.

// ------------- ErrInvalidData: некорректная модель на входе
// Некоторые ошибки GORM появляются ещё до обращения к базе —
// ORM проверяет, что вы вообще дали ей что-то разумное.

func DemoErrInvalidData(db *gorm.DB) error {
	// 1. Создание nil
	var u *models.User = nil
	res := db.Create(u)
	if errors.Is(res.Error, gorm.ErrInvalidData) {
		return fmt.Errorf("cannot create nil user: %w", res.Error)
	}

	// 2. Обновление «пустой» модели без ключа
	u2 := models.User{} // ID == 0
	res = db.Model(&u2).Updates(models.User{Name: "X"})
	if errors.Is(res.Error, gorm.ErrInvalidData) {
		// GORM не знает, какое условие WHERE строить
		return fmt.Errorf("missing primary key for update: %w", res.Error)
	}

	return nil
}

// ************************* Проверка входных данных до запроса *****************************
// Базу лучше воспринимать как «последнюю линию обороны».
// Всё, что можно проверить до неё, нужно проверять до неё: обязательные поля,
// формат email, длина пароля, диапазоны чисел, бизнес-правила уровня процесса.

// Паттерн удобно строить вокруг DTO: отдельных структур для входных данных
// с методами нормализации и валидации.

// ------------- DTO: нормализуем и валидируем запрос---------------------
// Простейший пример — регистрация пользователя:
type CreateUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (in *CreateUserDTO) Normalize() {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
}

func (in CreateUserDTO) Validate() error {
	// простейшие проверки для примера:
	if in.Email == "" || !strings.Contains(in.Email, "@") {
		return fmt.Errorf("invalid email: %w", gorm.ErrInvalidData) // ошибки маппятся на ErrInvalidData, чтобы контроллер мог отдать 400/422
	}
	if len(in.Password) < 8 {
		return fmt.Errorf("password too short: %w", gorm.ErrInvalidData) // ошибки маппятся на ErrInvalidData, чтобы контроллер мог отдать 400/422

	}
	return nil
}

// Сервисный слой:
func CreateUser(db *gorm.DB, in CreateUserDTO) (*models.User, error) {
	in.Normalize()

	if err := in.Validate(); err != nil {
		// можно сразу вернуть 400/422 без обращения к БД
		return nil, err
	}

	// сама модель (User) может дополнительно проверять себя в хуках (например, хэшировать пароль).
	u := models.User{
		Email:    in.Email,
		Password: in.Password, // хэшируем в хуке BeforeCreate
	}

	if err := db.Create(&u).Error; err != nil {
		return nil, fmt.Errorf("create user failed: %w", err)
	}

	return &u, nil
}

// ------------ Числовые диапазоны и бизнес-инварианты
// Более содержательный пример — создание заказа:
type ItemDTO struct {
	SKU   string  `json:"sku"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
}

// Здесь на уровне DTO сразу отсеиваются:
// пустой покупатель;
// пустой список товаров;
// нулевое или отрицательное количество;
// отрицательная цена.
func (i ItemDTO) Validate() error {
	if i.SKU == "" {
		return fmt.Errorf("sku required: %w", gorm.ErrInvalidData)
	}
	if i.Qty <= 0 {
		return fmt.Errorf("qty must be > 0: %w", gorm.ErrInvalidData)
	}
	if i.Price < 0 {
		return fmt.Errorf("price must be >= 0: %w", gorm.ErrInvalidData)
	}
	return nil
}

type CreateOrderDTO struct {
	CustomerID uint      `json:"customer_id"`
	Items      []ItemDTO `json:"items"`
}

func (o CreateOrderDTO) Validate() error {
	if o.CustomerID == 0 {
		return fmt.Errorf("customer_id required: %w", gorm.ErrInvalidData)
	}
	if len(o.Items) == 0 {
		return fmt.Errorf("order must contain items: %w", gorm.ErrInvalidData)
	}

	for i, it := range o.Items {
		if err := it.Validate(); err != nil {
			return fmt.Errorf("item %d: %w", i, err)
		}
	}

	return nil
}

// После такой проверки до базы доходят только осмысленные заказы.
// При этом ключевые правила всё равно дублируют в схеме (NOT NULL, CHECK, UNIQUE),
// чтобы защититься от гонок и «обхода» правил другим кодом.

// ------------ Простейшая защита от мусора: проверка ID
// Даже для простого GET /users/{id} не стоит сразу бежать в базу.
// Нулевой или отрицательный ID можно отбрасывать на входе:
func GetUser(db *gorm.DB, id uint) (*models.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("id must be positive: %w", gorm.ErrInvalidData)
	}

	var u models.User
	if err := db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("db error: %w", err)
	}

	return &u, nil
}

// такой код
// сразу отбрасывает заведомо некорректный запрос;
// различает «не найдено» и «сломалось».

// ================== Обработка ошибок миграций и транзакций ===================================
// Миграции и транзакции — точки, где ошибка особенно опасна:
// можно оставить схему в полуприменённом состоянии или часть данных обновлённой,
// а часть — нет. Задача кода — либо довести операцию до конца,
// либо откатить всё, не оставляя «пола на две плитки».

// ************** Ошибки миграций: если не получилось — падаем *************************
// AutoMigrate() и Migrator() работают поверх DDL-запросов (CREATE TABLE, ALTER TABLE и т. д.).
// Любая ошибка здесь должна остановить сервис: лучше не стартовать, чем работать на сломанной схеме.
func DemoAutoMigrateErr(db *gorm.DB) error {
	// Минимальный шаблон:
	if err := db.AutoMigrate(&models.User{}, &models.Order{}); err != nil {
		return fmt.Errorf("auto-migrate failed: %w", err)
	}
	return nil
}

// Точечные изменения через Migrator():
func DemoMigrator(db *gorm.DB) error {
	m := db.Migrator()

	if !m.HasTable(&models.User{}) {
		if err := m.CreateTable(&models.User{}); err != nil {
			return fmt.Errorf("create table users: %w", err)
		}
	}

	if !m.HasColumn(&models.User{}, "Phone") {
		if err := m.AddColumn(&models.User{}, "Phone"); err != nil {
			return fmt.Errorf("add column users.phone: %w", err)
		}
	}
	return nil
}

// Любая ошибка оборачивается и пробрасывается выше. Приложение не продолжает работу «как ни в чём не бывало».

// ----------------- Миграции данных в транзакции ----------------
// Если нужно не только изменить схему, но и подготовить данные,
// всё это заворачивается в одну транзакцию:
func DemoMigrationInTransaction(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		// нормализуем старые данные
		if err := tx.Exec(`
		UPDATE orders
		   SET total = 0
		 WHERE total IS NULL
	`).Error; err != nil {
			return fmt.Errorf("normalize totals: %w", err)
		}

		// меняем схему
		if err := tx.Exec(`
		ALTER TABLE orders
		ALTER COLUMN total SET NOT NULL
	`).Error; err != nil {
			return fmt.Errorf("set total NOT NULL: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("data migration failed: %w", err)
	}
	return nil
}

// Если любой шаг внутри блока вернёт ошибку,
// GORM сделает ROLLBACK, и база останется в исходном состоянии.

// ----------- Ошибки в транзакциях: различаем причины -----------
// Типичный паттерн работы с транзакцией:
func Withdraw(db *gorm.DB, id uint, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be > 0: %w", gorm.ErrInvalidData)
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		var u models.User

		// читаем с блокировкой на обновление
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&u, id).Error; err != nil {
			return err // может быть ErrRecordNotFound или ошибка драйвера
		}

		if u.Balance < amount {
			return fmt.Errorf("insufficient funds: %w", gorm.ErrInvalidData)
		}

		return tx.Model(&u).
			Update("balance", gorm.Expr("balance - ?", amount)).Error
	})

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// пользователь не найден — 404
		return err
	case errors.Is(err, gorm.ErrInvalidData):
		// некорректный ввод / бизнес-условие — 422
		return err
	case err != nil:
		// инфраструктурная проблема — 500
		return fmt.Errorf("withdraw tx failed: %w", err)
	default:
		return nil
	}
}

// Транзакция:

// либо успешно завершится (COMMIT);
// либо вернёт ошибку, и GORM сделает ROLLBACK.
// Важно, что вы чётко различаете:

// «нет записи» (ErrRecordNotFound);
// «условия операции не выполнены» (ErrInvalidData);
// «что-то сломалось» (драйвер, тайм-аут, дедлок и т. д.).

// ------------ Контекст и тайм-ауты ----------------
// Ещё одна типовая причина падения транзакций — запрос просто висит слишком долго.
// Правильно ограничивать время выполнения через context:
func DemoTransactionContext(db *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// потенциально долгие действия
		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

// Если контекст истечёт:
// драйвер вернёт ошибку;
// GORM откатит транзакцию;
// соединение освободится.
// Так вы не держите блокировки и пул соединений бесконечно.

// ------------ «Дефер-откат» и паники ----------------
// При ручном управлении транзакцией шаблон почти всегда один:
func DemoTransactionDefer(db *gorm.DB) error {
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("begin tx failed: %w", tx.Error)
	}

	defer func() {
		_ = tx.Rollback() // безопасно вызвать даже после успешного Commit
	}()

	if err := tx.Create(&models.User{Name: "Анна"}).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}
	return nil
}

// Даже если вы выйдете из функции раньше или словите панику,
// Rollback() из defer снимет блокировки и вернёт соединение в пул.
