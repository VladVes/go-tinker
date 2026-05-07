package main

import (
	"context"
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

	if err := godotenv.Load("./cmd/lessons/orm-relations-download/.env"); err != nil {
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
	err = db.AutoMigrate(&models.User{}, &models.Order{})
	if err != nil {
		log.Fatalf("problem with automigrate profile model: %v ", err)
	}

	// -------------

	// ********************************** Transactions *********************************
	// транзакции: набор операций, который либо выполняется целиком, либо откатывается целиком.

	// ========================== db.Transaction() ==================
	// Вы передаёте функцию в вызов db.Transaction,
	// GORM сам открывает транзакцию, передаёт в неё tx *gorm.DB и ждёт результат.
	// Вернули nil — делается COMMIT.
	// Вернули ошибку — ROLLBACK.
	// При панике GORM тоже откатит транзакцию.
	err = db.Transaction(func(tx *gorm.DB) error {
		user := models.User{Name: "Janna", Email: "janna@mail.com", Age: 24}
		if err := tx.Create(&user).Error; err != nil {
			return err // любая ошибка — сигнал на откат
		}

		order := models.Order{UserID: user.ID, Number: "d55", Total: 1000}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// всё прошло без ошибок — транзакция закоммитится
		return nil
	})

	if err != nil {
		log.Println("транзакция откатилавсь: ", err)
	}
	// ---------- ещё пример:
	// Важно не забывать внутри транзакции везде использовать tx, а не исходный db. Любой вызов через db пойдёт мимо транзакции.

	err = db.Transaction(func(tx *gorm.DB) error {
		// здесь важно передать именно tx, а не db
		if err := CreateOrder(tx, 1, 2000); err != nil {
			return err
		}
		return nil
	})

	// ******************* WithContext() ****************************
	// Через WithContext() транзакцию можно привязать к контексту с тайм-аутом.
	// Если время вышло, запросы отменятся, транзакция откатится:
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// какие-то тяжёлые операции
		return nil
	})

	// ================ Begin, Commit, Rollback - ручное управление =====================
	// Иногда одной функции db.Transaction() недостаточно.
	// Например, границы транзакции шире одной функции,
	// или момент COMMIT должен зависеть от нескольких независимых веток логики.
	// В таких случаях можно управлять транзакцией вручную.

	// Транзакция начинается с db.Begin(). Этот вызов возвращает новый объект *gorm.DB,
	// привязанный к активной транзакции.
	// Далее все операции выполняются через него.
	// В конце — Commit() или Rollback().

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("не удалось открыть транзакцию: %v", tx.Error)
	}
	// Чтобы не размазывать Rollback() по коду, удобно повесить его в defer:
	defer tx.Rollback()
	// Если Commit() отработал успешно, повторный Rollback() ничего не испортит,
	// но гарантирует, что при раннем выходе или панике транзакция не «зависнет».

	user := models.User{Name: "Иван"}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return
	}

	order := models.Order{UserID: user.ID, Total: 1200}
	if err := tx.Create(&order).Error; err != nil {
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("ошибка коммита:", err)
	}
}

func CreateOrder(tx *gorm.DB, userID uint, total int) error {
	return tx.Create(&models.Order{UserID: userID, Total: total}).Error
}
