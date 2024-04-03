package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/golfz/fun-exercise-api/wallet"
	"log"
	"time"
)

type Wallet struct {
	ID         int       `postgres:"id"`
	UserID     int       `postgres:"user_id"`
	UserName   string    `postgres:"user_name"`
	WalletName string    `postgres:"wallet_name"`
	WalletType string    `postgres:"wallet_type"`
	Balance    float64   `postgres:"balance"`
	CreatedAt  time.Time `postgres:"created_at"`
}

func (p *Postgres) Wallets(filter wallet.Wallet) ([]wallet.Wallet, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	selectQuery := psql.Select("*").From("user_wallet")

	// prepare filter
	if filter.WalletType != "" {
		if !wallet.IsWalletTypeValid(filter.WalletType) {
			return []wallet.Wallet{}, nil
		}
		selectQuery = selectQuery.Where(sq.Eq{"wallet_type": filter.WalletType})
	}
	if filter.UserID != 0 {
		selectQuery = selectQuery.Where(sq.Eq{"user_id": filter.UserID})
	}

	sql, args, err := selectQuery.ToSql()
	log.Println(sql)
	if err != nil {
		return nil, err
	}

	rows, err := p.Db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []wallet.Wallet
	for rows.Next() {
		var w Wallet
		err := rows.Scan(&w.ID,
			&w.UserID, &w.UserName,
			&w.WalletName, &w.WalletType,
			&w.Balance, &w.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet.Wallet{
			ID:         w.ID,
			UserID:     w.UserID,
			UserName:   w.UserName,
			WalletName: w.WalletName,
			WalletType: w.WalletType,
			Balance:    w.Balance,
			CreatedAt:  w.CreatedAt,
		})
	}
	return wallets, nil
}
