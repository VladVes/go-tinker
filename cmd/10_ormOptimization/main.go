package main

import "time"

//  Оптимизация всегда про три вещи:
// как устроена схема (индексы и ограничения);
// как мы загружаем связи;
// сколько данных реально вытаскиваем из базы.

// ******************* Индексы и ограничения ************************
// Индекс — это структура данных (упрощённо, отсортированное дерево),
// по которой база быстро находит строки для WHERE, JOIN, ORDER BY и GROUP BY.
// Без индекса:
// SELECT * FROM users
// WHERE email = 'a@example.com';
// вынуждает СУБД просканировать всю таблицу.
// С индексом по email — это поиск в дереве, который практически не зависит от количества строк.

// Ограничения — это правила целостности: NOT NULL, UNIQUE, CHECK, FOREIGN KEY.
// Они не ускоряют запросы напрямую, но гарантируют, что в таблицу
// не попадут бессмысленные значения (например, отрицательная сумма или дублирующийся email).

// В GORM всё это задаётся прямо в модели.
// Пример: пользователь с уникальным email и индексом по дате создания:

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex;size:120;not null"`
	Name      string    `gorm:"index;size:80"`
	CreatedAt time.Time `gorm:"index"`
}

// uniqueIndex — создаст уникальный индекс по Email (нельзя вставить двух пользователей с одним email);
// index на Name и CreatedAt — обычные индексы для частых фильтров и сортировок.
// При миграции GORM сгенерирует команды вида:
// CREATE UNIQUE INDEX idx_users_email ON users (email);
// CREATE INDEX idx_users_name ON users (name);
// CREATE INDEX idx_users_created_at ON users (created_at);

// --------- Составной индекс ----
// Если вы часто фильтруете по сочетанию полей, нужен составной индекс:

type Post struct {
	ID        uint
	UserID    uint   `gorm:"index:idx_user_title,unique"`
	Title     string `gorm:"index:idx_user_title,unique"`
	CreatedAt time.Time
}

// Этот тег index:idx_user_title,unique говорит: создать один составной уникальный индекс по (user_id, title). Такой индекс ускоряет запросы:
// db.Where("user_id = ? AND title = ?", 42, "Hello").First(&post)
// и одновременно не даст одному пользователю создать два поста с одинаковым заголовком.

// ------ Органичения CHECK и NOT NULL ----
// Ограничения CHECK и NOT NULL описываются так же:
type Order struct {
	ID        uint
	Amount    float64 `gorm:"check:amount_non_negative,amount >= 0"`
	UserID    uint    `gorm:"not null"`
	CreatedAt time.Time
}

// При миграции появится ограничение:
// ALTER TABLE orders
// ADD constraint amount_non_negative CHECK (amount >= 0);
// Тонкая настройка возможна через Migrator:
// m := db.Migrator()
// if !m.HasIndex(&User{}, "idx_users_email") {
// 	_ = m.CreateIndex(&User{}, "Email")
// }

// Так можно постепенно добавлять индексы в существующий проект, ориентируясь на реальные запросы.
// Главная мысль: индексы ставят под конкретные паттерны использования
// (частые WHERE, JOIN, ORDER BY), а ограничения — под критичные бизнес-правила.
// Всё остальное — лишний груз:
// каждый индекс ускоряет чтение, но замедляет вставки и обновления,
// потому что базе нужно обновлять индексы при каждой операции.

// ************ Избежание N+1 при загрузке связей ********************************
// Классическая проблема ORM — запрос N+1. Сценарий выглядит так:

// - Вы делаете один запрос за списком сущностей (например, пользователей).
// - Потом в цикле для каждого пользователя отдельно вытаскиваете его посты.
// В результате база получает 1 запрос на список + N запросов на связи.
// На десятках пользователей это быстро превращается в лавину.

// Плохой код:
// var users []User
// db.Find(&users) // SELECT * FROM users;

// for i := range users {
// 	// для каждого пользователя отдельный запрос
// 	if err := db.Model(&users[i]).Association("Posts").Find(&users[i].Posts); err != nil {
// 		// ...
// 	}
// }
//Логи будут такими:
// SELECT * FROM users;
// SELECT * FROM posts WHERE user_id = 1;
// SELECT * FROM posts WHERE user_id = 2;
// SELECT * FROM posts WHERE user_id = 3;
// ...

// Правильный подход — жадная загрузка (eager loading) через Preload():
// var users []User
// db.Preload("Posts").Find(&users)
// GORM сделает всего два запроса:
// SELECT * FROM users;
// SELECT * FROM posts WHERE user_id IN (1,2,3,...);
// и разложит посты по пользователям в памяти.

// Можно сразу фильтровать связи:
// db.Preload("Posts", "published = ?", true).Find(&users)
// SQL получится:
// SELECT * FROM posts
// WHERE user_id IN (1,2,3,...) AND published = true;
// Если у вас вложенные связи, их можно подгружать цепочкой:
// db.Preload("Posts.Comments").Preload("Profile").Find(&users)

// Иногда нужно не просто подгрузить связи, но и отфильтровать сущности по ним.
// Тогда помогают Joins():
// var users []User
// db.Joins("JOIN posts ON posts.user_id = users.id").
//    Where("posts.likes > ?", 100).
//    Find(&users)

// Что важно:
// - N+1 — это не баг GORM, а побочный эффект ленивых запросов к связям.
// - Чтобы его избежать, нужно осознанно использовать Preload() и/или Joins(),
// а не ходить к association в цикле.
// - Preload("*") — тяжёлая артиллерия. Она убирает N+1, но подгружает всё подряд,
// что может стать отдельной проблемой производительности.

// ****************** Выборка только нужных полей ****************************
// Даже когда запросов немного и индексы настроены, можно легко «убить» производительность,
// таща лишние данные. GORM по умолчанию генерирует SELECT *, то есть выбирает
// все колонки модели. Для сущностей с большим количеством полей, JSON-колонками
// или BLOB’ами это лишняя нагрузка на базу, сеть и парсер драйвера.
// Метод Select() позволяет задать, какие именно колонки нужны.
// Простейший пример:

// var users []User
// db.Select("id", "email", "created_at").Find(&users)

// SQL:

// SELECT
//     id,
//     email,
//     created_at
// FROM users;

// ------- Scan() -----
// Отдельный часто используемый паттерн — выборка в «облегчённую» структуру через Scan():
// type UserLite struct {
// 	ID    uint
// 	Email string
// }
// var result []UserLite
// db.Model(&User{}).
//    Select("id", "email").
//    Scan(&result)
// Это удобно для API-слоя: одна модель (User) для полной работы с таблицей,
// другая (UserLite) — для DTO в ответе.

// Select() хорошо сочетается с Preload():
// var users []User
// db.Select("id", "name").
//    Preload("Profile", func(db *gorm.DB) *gorm.DB {
// 	   return db.Select("user_id", "bio")
//    }).
//    Find(&users)
// SQL будет примерно таким:

// SELECT id, name FROM users;
// SELECT user_id, bio FROM profiles WHERE user_id IN (...);

// Производительность кода с GORM держится на трёх опорах:

// - Схема: индексы под реальные запросы и ограничения на критичные инварианты.
// Когда часто фильтруете или сортируете по полю
// - Запросы к связям: явное использование Preload/Joins, чтобы не попадать в N+1
// и не тянуть лишнее.
// - Объём данных: Select() и лёгкие DTO вместо бесконечных SELECT *.
