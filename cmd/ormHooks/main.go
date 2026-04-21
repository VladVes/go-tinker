package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
}

// модели могут быть «умнее»: проверять себя перед сохранением,
// автоматически считать derived-поля, вести аудит, запрещать опасные операции.
// Всё это делается через хуки — специальные методы,
// которые GORM вызывает в ключевые моменты жизненного цикла записи.

// Хуки позволяют перенести часть инвариантов и бизнес-правил ближе к данным.
// Модель перестаёт быть «тупым мешком полей» и начинает сама следить за своей целостностью:
// нормализует email, хэширует пароль, не даёт удалить оплачённый заказ и т. д.

// GORM ищет на типе модели методы с особыми именами и вызывает их автоматически.
// Каждый такой метод получает *gorm.DB, работает в той же транзакции,
// что и основная операция, и может вернуть ошибку, чтобы эту операцию отменить.

// Основные хуки:
// 	BeforeCreate() / AfterCreate() — до и после вставки новой записи;
// BeforeSave() / AfterSave() — до и после любого сохранения (и создания, и обновления);
// BeforeUpdate() / AfterUpdate() — до и после обновления;
// BeforeDelete() / AfterDelete() — до и после удаления;
// AfterFind() — сразу после выборки записи из базы.

// Хук всегда объявляется на указателе:
// func (u *User) BeforeCreate(tx *gorm.DB) error {
// 	// ...
// 	return nil
// }

// Если метод вернул error, GORM:
// прерывает текущую операцию (Create(), Save(), Delete() и т. д.);
// откатывает транзакцию, если операция была внутри транзакции;
// возвращает ошибку вызывающему коду.
// Это удобно: вся критичная валидация и подготовка данных может жить прямо в модели, а не размазываться по сервисам.

// ******************** Изменение поведения операций через хуки **********************
// Хуки — не только про проверки. Они позволяют менять поведение операций:
// автоматически дополнять запись, вычислять derived-поля, писать аудит,
// подготавливать данные для чтения.

// ------- Пример: продукт, который сам считает цену с НДС перед любым сохранением:
type Product struct {
	ID           uint
	Name         string
	Price        int
	PriceWithVAT int
	UpdatedAt    time.Time
}

func (p *Product) BeforeSave(tx *gorm.DB) error {
	if p.Price < 0 {
		return fmt.Errorf("price must be non-negative")
	}

	// 20% НДС
	p.PriceWithVAT = int(math.Round(float64(p.Price) * 1.2))
	return nil
}

// Как только где-то в коде выполнится: db.Save(&product)
// GORM сначала вызовет BeforeSave(), модель проверит цену и пересчитает PriceWithVAT,
// и только затем ORM отправит UPDATE в базу.
// Если цена некорректная, BeforeSave() вернёт ошибку, запрос не уйдёт в базу, а транзакция откатится.

// -------- Пример: обогащение данных после чтения через AfterFind().
// Пусть у пользователя в профиле хранится часовой пояс, а в коде нам удобно иметь
// «текущее локальное время» в этом поясе, но хранить его в базе бессмысленно.
type Profile struct {
	ID       uint
	Name     string
	Timezone string
	LocalNow time.Time `gorm:"-"` // не хранится в таблице
}

func (p *Profile) AfterFind(tx *gorm.DB) error {
	loc, err := time.LoadLocation(p.Timezone)
	if err != nil {
		loc = time.UTC
	}
	p.LocalNow = time.Now().In(loc)
	return nil
}

// При выборке: db.First(&profile, 1)
// GORM выполнит SELECT по таблице profiles,
// распакует поля в структуру,
// вызовет AfterFind(), который заполнит LocalNow.
// В итоге в коде сразу доступна удобная, «готовая к использованию» структура, а в базе по-прежнему хранятся только необходимые поля.

// -------- Пример: запрет опасного удаления через BeforeDelete().
// Заказ со статусом paid вы хотите только отменять или рефандить, но не удалять физически:
type Order struct {
	ID        uint
	Status    string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (o *Order) BeforeDelete(tx *gorm.DB) error {
	if o.Status == "paid" {
		return fmt.Errorf("paid orders cannot be deleted")
	}
	return nil
}

// Вызов: db.Delete(&order)
// для оплаченного заказа закончится ошибкой, а строка в базе останется на месте.
// Правило живёт рядом с моделью и сработает независимо от того, из какого участка кода попытались удалить заказ.

// ------ Пример: аудит
// После создания пользователя можно записать событие в отдельную таблицу в рамках той же транзакции:
type AuditLog struct {
	ID        uint
	Entity    string
	EntityID  uint
	Action    string
	CreatedAt time.Time
}

type Usr struct {
	ID    uint
	Name  string
	Email string
}

func (u *Usr) AfterCreate(tx *gorm.DB) error {
	log := AuditLog{
		Entity:   "users",
		EntityID: u.ID,
		Action:   "create",
	}
	return tx.Create(&log).Error
}

// Если по какой-то причине запись аудита не сохранилась,
// AfterCreate() вернёт ошибку, и вся операция создания пользователя откатится.

// -------------- Автоматическое хэширование пароля ---------------------
// Правильный вариант — всегда сохранять только криптографический хэш (bcrypt, Argon2, scrypt).
// В GORM этим удобно заниматься на уровне модели.
type User struct {
	ID        uint
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Задача: при любой вставке:
// нормализовать email (обрезать пробелы, привести к нижнему регистру);
// проверить длину пароля;
// заменить пароль на bcrypt-хэш.

const MinPasswordLength = 8

func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Нормализуем email.
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	// Проверяем минимальную длину пароля.
	if len(u.Password) < MinPasswordLength {
		return fmt.Errorf("password too short: minimum %d characters", MinPasswordLength)
	}
	// Генерируем bcrypt-хэш.
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	u.Password = string(hashed)
	return nil
}

// Теперь любая попытка создать пользователя:
// err := db.Create(&User{
// 	Email:    "User@Mail.Com ",
// 	Password: "secret123",
// }).Error
// пройдёт через BeforeCreate(). В базу уйдёт строка вроде:
// email:    "user@mail.com"
// password: "$2a$10$ZJBlE....2aOSQmv9T.N."

// Осталось покрыть обновление пароля. Для этого используют BeforeSave() и проверку изменённых полей:
func (u *User) BeforeSave(tx *gorm.DB) error {
	// email тоже нормализуем при любом сохранении
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	// если пароль не менялся — ничего не делаем
	if !tx.Statement.Changed("Password") {
		return nil
	}

	if len(u.Password) < 8 {
		return fmt.Errorf("password too short")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// Теперь при любой операции вроде:
// db.Model(&user).Updates(map[string]interface{}{
// 	"Password": "newSecret123",
// })
// GORM Увидит изменение поля Password.
// Вызовет BeforeSave().
// Захэширует новый пароль.
// Сохранит в базе только хэш.
// Если пароль не менялся, tx.Statement.Changed("Password") вернёт false, и хук не будет трогать существующее значение.

// ********************************* Практические нюансы хуков
// Несколько важных моментов, о которых стоит помнить:

// Хуки работают внутри той же транзакции. Если вы вызываете db.Transaction(),
// все Before*/After* выполняются в рамках одного BEGIN/COMMIT.

// Любая ошибка из хука откатывает операцию. Это удобно для валидации,
// но опасно, если держать внутри долгие операции.

// Тяжёлую работу (HTTP-запросы, запись в внешние очереди) лучше не делать непосредственно в хуках.
// Правильный паттерн — писать событие в outbox-таблицу в AfterCreate() и обрабатывать его асинхронным воркером.

// Если по какой-то причине нужно явно отключить хуки для конкретного запроса,
// можно использовать сессию: db.Session(&gorm.Session{SkipHooks: true}).Create(&user),
// но это редкий сценарий; по умолчанию лучше считать, что хуки всегда включены.
