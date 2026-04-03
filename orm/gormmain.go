package orm

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DSN = "host=localhost user=postgres password=mysecretpassword dbname=postgres port=5432 sslmode=disable"

type User struct {
	ID   uint   // первичный ключ
	Name string // имя пользователя
}

func GormRun() {

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		// логирование ошибки подключения и завершение программы
		log.Fatalf("DB connection error: %v", err)
	}
	log.Println("Соединение с базой установлено")

	// AutoMigrate создаёт таблицу users, если её ещё нет,
	// и обновляет схему при изменении структуры
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("ошибка миграции схемы: %v", err)
	}
	log.Println("Таблица users готова к работе")

}
