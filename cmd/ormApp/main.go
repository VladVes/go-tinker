package main

import (
	"errors"
	"fmt"
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

// ---------------------------------------------- MODELS -----------------------------------
// Точнее управлять отображением полей в таблицу позволяют теги GORM.
// Эти теги описываются в обратных кавычках после объявления поля и задают имя столбца,
// тип, ограничения и ключи. Пример модели с настроенными тегами показывает,
// как структура управляет схемой базы:
type User struct {
	// первичный ключ Использование uint для ID — общепринятая практика в Go‑проектах с GOR
	// uint - беззнаковый целочисленный универсальный для 32 и 64 битных систем
	// если поле имеет имя ID и тип uint gorm понимает это как первичный ключ
	ID    uint   `gorm:"primaryKey"`                    // Явное указание на первичный ключ.
	Name  string `gorm:"column:login;size:50;not null"` // Имя столбца login, длина 50, уникальное и не NULL. //  директиву unique убрал что бы не чистить таблицу перед каждым запуском примера
	Email string `gorm:"type:varchar(100)"`             // Явный тип столбца для email.
	Age   int    `gorm:"default:18"`                    // Значение по умолчанию.
	// CreatedAt, UpdatedAt и DeletedAt. Большинство таблиц нуждается в информации о том, когда запись создана,
	// когда изменена и была ли удалена логически.
	// GORM умеет управлять такими полями автоматически
	CreatedAt time.Time // Время создания записи
	UpdatedAt time.Time // Время последнего обновления.

	// При удалении через db.Delete запись не исчезает физически:
	// GORM устанавливает в DeletedAt текущий момент времени, а при последующих
	// выборках с обычными методами такие строки больше не попадают в результаты.
	// Чтобы увидеть все записи, включая помеченные как удалённые, используется
	// Unscoped(). Для физического удаления строк с ненужными данными также
	// применяется Unscoped() вместе с Delete().
	DeletedAt gorm.DeletedAt `gorm:"index"` // Пометка времени удаления и индекс по этому полю.
}

// Связи между таблицами также управляются через теги.
// Внешний ключ и поведение при обновлении и удалении можно описать в самой модели.
type Order struct {
	ID uint `gorm:"primaryKey"`
	// Поле UserID хранит идентификатор пользователя и получает ограничения NOT NULL и INDEX.
	UserID uint `gorm:"not null;index"` // Поле внешнего ключа, индекс и запрет NULL.
	// Поле User задаёт связь: foreignKey указывает, что для связи используется UserID,
	// параметр references объявляет, на какое поле в структуре User идёт ссылка.
	// А constraint описывает поведение СУБД: при обновлении ключа в таблице пользователей заказы обновятся каскадно,
	// а при удалении пользователя ссылки в заказах превратятся в NULL.
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// ------------------ gorm.Model -------
// Для удобства GORM предоставляет встроенный тип gorm.Model.
// Этот тип уже включает в себя стандартный набор полей:
// ID, CreatedAt, UpdatedAt и DeletedAt.
// Любая модель, встраивающая gorm.Model, автоматически получает эти поля и стандартное поведение:
type Product struct {
	gorm.Model // ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string
	Price      int
}

// --------------- Методы для изменеия дефолтного создания таблиц gorm -----------------------
// Имя таблицы по умолчанию удобно не всегда.
// Когда в существующей базе используются нестандартные наименования,
// GORM нужно подсказать правильное имя. Для этого структура может объявить метод TableName:

// т.к. метод не работает с полями структуры можно его определить
// без создания переменной для получаетля значения т.е. это аналог static методов класса в других яп
// но обращаться к методу нужно будет всё равно использую экземпляр tName := usr.TableName()

// func (User) TableName() string {
// Теперь все операции с моделью User будут обращаться к таблице app_users.
// Такой приём помогает адаптировать ORM к уже существующей схеме без её переписывания.
// 	return "app_users"
// }

func main() {
	// ------------------------------- Env var -------------------------------------------------------------------
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	fmt.Printf("Data Source Name из env var: DATABASE_URL = %s\n", dsn)
	if dsn == "" {
		dsn = defaultDsn
	}

	// ------------------------------- logger -----------------------------------------------------------------

	// Создание нового логгера с настройками
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // базовый вывод в консоль
		logger.Config{
			SlowThreshold: time.Second, // порог для медленных запросов
			LogLevel:      logger.Info, // подробный уровень логирования
			Colorful:      true,        // цветной вывод для удобства
		},
	)

	// ------------------------------- Подключение и конфигурация -----------------------------------------------------------------

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// В таком режиме при каждом запросе GORM будет печатать текст SQL,
		// количество затронутых строк и время выполнения. Это помогает увидеть реальную картину нагрузки,
		// найти медленные места и понять, как ORM строит запросы под капотом.
		Logger: newLogger,

		// отключить обёртку каждой операции в транзакцию и включить кэш подготовленных выражений:
		SkipDefaultTransaction: true, // ускоряет типовые операции, если транзакции контролируются на более высоком уровне
		PrepareStmt:            true, //  полезен в сервисах с повторяющимися запросами: подготовленные выражения переиспользуются, и база тратит меньше времени на разбор SQL

		// Глобальное управление именами таблиц возможно через стратегию именования.
		// Если все таблицы должны иметь общий префикс, его удобно задать в конфигурации GORM:
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "app_",
			// После этого User связан с app_users, Order() — с app_orders и так далее.
			// Такой подход выравнивает схему по единым правилам и помогает избежать конфликтов имён.
		},
	})
	if err != nil {
		// логирование ошибки подключения и завершение программы
		log.Fatalf("DB connection error: %v", err)
	}
	log.Println("Соединение с базой установлено")

	// ------------------------------- Пул соединений db.DB и настройка-----------------------------------------------------

	// важную роль играет настройка пула соединений.
	// GORM использует стандартный пакет database/sql и создаёт пул подключений к базе.
	// Управление этим пулом выполняется через метод DB() объекта *gorm.DB:
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("ошибка доступа к пулу соединений: %v", err)
	}
	// Ping проверяет, что соединение живое и база отвечает
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("ошибка пинга базы: %v", err)
	}

	// Максимальное число открытых соединений к базе
	sqlDB.SetMaxOpenConns(10)
	// Максимальное время жизни соединения
	sqlDB.SetConnMaxLifetime(time.Hour)
	log.Println("Пул соединений настроен и готов к работе")
	// Такая настройка помогает удерживать баланс между производительностью и нагрузкой на базу.
	// Пул не открывает бесконечное количество соединений и не держит их пожизненно, а работает по заданным ограничениям.
	// В продакшене это особенно важно: правильные значения параметров позволяют выдерживать нагрузку и не исчерпывать ресурсы СУБД.

	// ------------------------------- DryRun ----------------------------------------------------------------------------------

	// Для проверки того, какой запрос GORM сформирует без реального выполнения, можно использовать сессию с параметром DryRun().
	// В этом режиме ORM строит SQL, но не отправляет его в базу

	// DryRun подходит для отладки сложных цепочек запросов и помогает убедиться,
	// что GORM строит именно тот SQL, который ожидается, ещё до взаимодействия с реальной базой.
	tx := db.Session(&gorm.Session{
		DryRun: true, // режим генерации SQL без выполнения
	})

	// Формирование SELECT-запроса без обращения к базе
	stmt := tx.First(&User{}, 1).Statement

	log.Println("DryRun SQL: ", stmt.SQL.String())

	// ------------------------------- Automigrate первый пример  -------------------------------------------------------

	// AutoMigrate создаёт таблицу users, если её ещё нет,
	// и обновляет схему при изменении структуры
	// При изменении структуры (например, при добавлении нового поля) GORM попытается аккуратно обновить схему.
	// Такой механизм удобен на ранних этапах разработки, когда модель данных ещё активно меняется.
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("ошибка миграции схемы: %v", err)
	}
	log.Println("Схема для users синхронизирована")

	// Важно, что AutoMigrate() ориентирован на безопасные изменения.
	// Он не удаляет столбцы, даже если поле исчезло из структуры.
	// Он не меняет типы колонок, если они уже существуют.
	// Такое поведение защищает от потери данных, но означает, что сложные изменения схемы
	// (переименование поля, сжатие типов, удаление колонок) требуют явных миграций на уровне SQL
	// или отдельных инструментов вроде goose.
	// GORM в этом смысле ведёт себя как мягкий помощник: добавляет недостающее, но не ломает существующее.

	// В простых проектах AutoMigrate() способен полностью закрыть потребности в эволюции схемы.
	// На этапе активной разработки он позволяет быстро изменять модели и получить актуальную структуру без написания DDL-скриптов.
	// В более зрелых системах его обычно используют осторожно: для небольших дополнений и лайтовых изменений,
	// а крупные рефакторинги схемы оформляют в виде явных миграций с контролем версий.

	// ******************************************************************************************************

	// ------------------------------- Simple Queries examples -----------------------------------------------------------------------------
	db.Create(&User{Name: "Alice", Email: "alice@mail.com", Age: 25})
	var usr1 User
	db.First(&usr1, 1)

	log.Println("Пользователь db.First: ", usr1)

	// ------------------------------- Queries examples with errors handling -----------------------------------------------------------------------------
	if err := db.Create(&User{Name: "Анна", Email: "anna@example.com"}).Error; err != nil {
		log.Fatalf("ошибка вставки: %v", err)
	}

	var usr2 User
	if err := db.First(&usr2).Error; err != nil {
		log.Fatalf("ошибка чтения: %v", err)
	}

	log.Printf("пользователь загружен: %s <%s>", usr2.Name, usr2.Email)

	// ------------------------------ Movies model example-------------

	var movie models.Movie

	if err := db.First(&movie).Error; err != nil {
		log.Fatalf("ошбика получения записи из movies: %v", err)
	}

	log.Printf("Первая запись из таблицы movies: %v", movie)

	// -----------------------После доработки Movies (добавление новых столбцов) -------------
	if err := db.AutoMigrate(&models.Movie{}); err != nil {
		log.Fatalf("ошибка миграции схемы: %v", err)
	}
	log.Println("Схема для movies синхронизирована")

	// *************************************************************************
	// ----------------------------- CRUD ------------------------------------
	// Когда программа вызывает
	// db.Create, библиотека разбирает структуру, читает теги и значения полей,
	// строит INSERT, подставляет параметры и отправляет запрос в базу.
	//
	// Методы First(), Find() и Where() превращаются в SELECT с нужными условиями
	// и ограничениями.
	//
	// Updates() и Update() становятся UPDATE,
	//
	// Delete() — либо DELETE, либо UPDATE с заполнением поля DeletedAt при мягком удалении

}

// ------------------------- CRUD - CREATE ---------------------------
func CreateUser(db *gorm.DB) error {
	// Формируется новая структура без ID — ключ задаст база.
	user := User{
		Name:  "Екатерина",
		Email: "katka@mail.com",
	}

	// Create строит INSERT и выполняет его.
	result := db.Create(&user)

	// Проверка ошибки вставки.
	if result.Error != nil {
		return result.Error
	}

	// После вставки база вернёт ID, GORM запишет его в user.ID.
	log.Println("Создан пользователь с ID:", user.ID)

	// Через RowsAffected можно узнать, сколько строк добавлено.
	log.Println("Добавлено строк:", result.RowsAffected)

	return nil
}

// --------------- CRUD - CREATE with slice ---------------------------
// При массовой вставке срез структур превращается в один батчевый запрос.
// Программа создаёт набор значений, а GORM формирует INSERT с несколькими
// строками
// Такой подход снижает количество запросов и ускоряет добавление большого количества данных

func CreateManyUsers(db *gorm.DB) error {
	users := []User{
		{Name: "Игорь", Email: "igor@mail.com"},
		{Name: "Мария", Email: "maria@mail.com"},
	}

	// Один вызов Create для среза создаст несколько строк за один SQL-запрос.
	if err := db.Create(&users).Error; err != nil {
		return err
	}

	// После вставки у каждого элемента среза появится свой ID.
	for _, u := range users {
		log.Println("Пользователь:", u.Name, "ID:", u.ID)
	}

	return nil
}

// --------------- CRUD - READ: First, Find, Where ---------------------------

// First
// Метод First() ориентирован на получение одной записи.
// Если передать только структуру, GORM выберет первую строку по возрастанию первичного ключа:
func FirstUser(db *gorm.DB) (User, error) {
	var user User

	// Без условий: первая запись в таблице users.
	result := db.First(&user)

	if result.Error != nil {
		return User{}, result.Error
	}

	log.Println("Первый пользователь:", user.ID, user.Name)
	return user, nil
}

// First with id
// Если программе известен первичный ключ, его можно передать вторым
// аргументом. В этом случае GORM сформирует запрос с
// WHERE id = ? и LIMIT 1:
func FindUserByID(db *gorm.DB, id uint) (User, error) {
	var user User

	// Выборка по первичному ключу.
	result := db.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("Пользователь с таким ID не найден")
		return User{}, result.Error
	}

	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}

// Find
// Для выборки нескольких строк используется Find().
// Этот метод заполняет срез структур и может комбинироваться с фильтрами
// и сортировкой:
func ListUsers(db *gorm.DB) ([]User, error) {
	var users []User

	// Получение всех записей без условий.
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	log.Println("Найдено пользователей:", len(users))
	return users, nil
}

// Where
// Фильтрация добавляется методом Where().
// Он может принимать SQL-условия с подстановками,
// структуру или карту. Пример выборки по возрасту и почте
// показывает работу с текстовым условием:
func FindAdultsWithMail(db *gorm.DB) ([]User, error) {
	var users []User

	// Условие age > 25 и домен почты *@mail.com.
	query := db.Where("age > ?", 25).Where("email LIKE ?", "%@mail.com")

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// Select
// Для выборки конкретных колонок используется Select().
// Тогда в структуру будут загружены только запрошенные поля,
// а остальные останутся нулевыми значениями:
func FindNamesAndEmails(db *gorm.DB) ([]User, error) {
	var users []User

	// Загрузка только name и email. Остальные поля не читаются из базы.
	if err := db.Select("name", "email").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// Во всех этих случаях GORM формирует объект Statement с именем таблицы,
// списком полей, фильтрами, сортировкой и ограничениями.
// Затем ORM строит SQL под конкретную СУБД, выполняет запрос и сканирует
// строки результата в структуры Go, автоматически приводя типы.

// --------------- CRUD - UPDATE: Save, Update, Updates ----------------

// Save
// Метод Save() воспринимает структуру как снимок всей записи.
// Если у объекта заполнен первичный ключ, GORM считает, что строка уже
// существует, и выполняет UPDATE. Если ключ пустой, ORM воспринимает
// структуру как новую запись и делает INSERT. Такой подход делает Save()
// универсальным, но иногда слишком широким: он обновляет все поля,
// включая нулевые.
func UpdateUserWithSave(db *gorm.DB, id uint) error {
	var user User

	// Сначала выбирается запись по ID.
	if err := db.First(&user, id).Error; err != nil {
		return err
	}

	// Изменение полей структуры в памяти.
	user.Name = "Анна Петрова"
	user.Email = "new@mail.com"

	// Save обновляет все поля записи в базе.
	if err := db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// Updates Update
// Метод Updates() ориентирован на частичные изменения.
// Он принимает либо структуру, либо карту и модифицирует
// только указанные поля. Остальные колонки остаются без изменений.
// При передаче структуры нулевые значения пропускаются,
// а при передаче map обновляются все перечисленные поля,
// даже если они равны нулю.
func PartialUpdateUser(db *gorm.DB, id uint) error {
	var user User

	if err := db.First(&user, id).Error; err != nil {
		return err
	}

	// Обновление нескольких полей через структуру:
	// нулевые значения будут проигнорированы.
	if err := db.Model(&user).Updates(User{
		Name: "Анна",
		Age:  30,
	}).Error; err != nil {
		return err
	}

	// Обновление конкретного поля через Update.
	if err := db.Model(&user).Update("Email", "updated@mail.com").Error; err != nil {
		return err
	}

	return nil
}

// Когда важно записать и нулевые значения, программа передаёт карту:
var ErrNotFound = errors.New("Error not found")

func ResetUserAge(db *gorm.DB, id uint) error {
	var user User

	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("find user %d: %w", id, ErrNotFound)
		}
		return fmt.Errorf("find user %d: %w", id, err)
	}

	data := map[string]any{"age": 0}

	if err := db.Model(&user).Updates(data).Error; err != nil {
		return fmt.Errorf("update user %d age: %w", id, err)
	}

	return nil
}

// --------------- CRUD - DELETE  ----------------
// В GORM удаление может быть двух видов: физическое удаление строки и мягкое
// удаление, когда запись остаётся в таблице, но помечается как скрытая
// с помощью DeletedAt.
// Если модель не содержит DeletedAt, Delete() превращается в прямой
// DELETE:
func RemoveUserHard(db *gorm.DB, id uint) error {
	// Удаление по ID без предварительной загрузки структуры.
	result := db.Delete(&User{}, id)

	if result.Error != nil {
		return result.Error
	}

	log.Println("Удалено строк:", result.RowsAffected)
	return nil
}

// Когда структура включает поле DeletedAt типа gorm.DeletedAt,
// GORM переключается на мягкое удаление. Вместо уничтожения строки
// ORM устанавливает в DeletedAt текущее время, а при обычных выборках
// такие записи автоматически игнорируются:

type SoftUser struct {
	ID        uint
	Name      string
	Email     string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func SoftDeleteUser(db *gorm.DB, id uint) error {
	var user SoftUser

	if err := db.First(&user, id).Error; err != nil {
		return err
	}

	// Мягкое удаление: обновление поля deleted_at.
	if err := db.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

// Unscoped
// При последующих вызовах Find() и First() GORM добавляет в запрос
// условие deleted_at IS NULL и исключает помеченные строки.
// Чтобы увидеть все записи, включая мягко удалённые,
// программа использует Unscoped().
// Тот же метод применяется для окончательного удаления:

func HardDeleteSoftUser(db *gorm.DB, id uint) error {
	var user SoftUser

	if err := db.Unscoped().First(&user, id).Error; err != nil {
		return err
	}

	// Полное удаление строки из таблицы.
	if err := .Error; err != nil {
		return err
	}

	return nil
}

// Такой механизм позволяет сначала безопасно скрывать данные,
// а затем, при необходимости, выполнять физическую очистку базы
// в отдельном процессе или по отдельному сценарию.
