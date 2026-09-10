// Реализуйте функцию ApplyCashbackRaw, которая начисляет кэшбек пользователям у которых есть покупки за период.

// Пример использования:

// from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
// to := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

// // начисляем кешбек в размере 10
// if err := repository.ApplyCashbackRaw(conn, from, to, 10); err != nil {
// 	log.Fatalf("cashback failed: %v", err)
// }
// В этом задании используется база данных SQLite, учитывайте это при написании запросов.

package repository

import (
	"time"

	"gorm.io/gorm"
)

// Customer represents a customer.
type Customer struct {
	ID        uint
	Name      string
	Wallet    *Wallet
	Purchases []Purchase
}

// Purchase represents a purchase made by a customer.
type Purchase struct {
	ID         uint
	CustomerID uint
	Amount     int64
	CreatedAt  time.Time
	Customer   *Customer `gorm:"foreignKey:CustomerID"`
}

// Wallet represents a customer wallet.
type Wallet struct {
	ID         uint
	CustomerID uint
	Balance    int64
	Customer   *Customer `gorm:"foreignKey:CustomerID"`
}

// ApplyCashbackRaw applies cashback for customers with purchases in period.
func ApplyCashbackRaw(conn *gorm.DB, from, to time.Time, amount int64) error {
	// BEGIN (write your solution here)
	query := `
		UPDATE wallets
		SET balance = balance + ?
		WHERE EXISTS (
			SELECT 1
			FROM purchases AS p
			WHERE
				p.customer_id = wallets.customer_id
				AND p.created_at >= ?
				AND p.created_at <= ?
		)
	`
	return conn.Exec(
		query,
		amount,
		from,
		to,
	).Error

	// END
}
