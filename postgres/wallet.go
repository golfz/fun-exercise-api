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

	wallets := make([]wallet.Wallet, 0)
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

func (p *Postgres) CreateWallet(wallet *wallet.Wallet) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	insertQuery := psql.Insert("user_wallet").
		Columns("user_id", "user_name", "wallet_name", "wallet_type", "balance").
		Values(wallet.UserID, wallet.UserName, wallet.WalletName, wallet.WalletType, wallet.Balance).
		Suffix("RETURNING id")

	sql, args, err := insertQuery.ToSql()
	if err != nil {
		return err
	}

	err = p.Db.QueryRow(sql, args...).Scan(&wallet.ID)
	if err != nil {
		return err
	}

	selectSQL := psql.Select("created_at").From("user_wallet").Where(sq.Eq{"id": wallet.ID})
	sql, args, err = selectSQL.ToSql()
	if err != nil {
		return err
	}

	err = p.Db.QueryRow(sql, args...).Scan(&wallet.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateWallet(wallet *wallet.Wallet) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updateQuery := psql.Update("user_wallet").
		Set("balance", wallet.Balance).
		Where(sq.Eq{"id": wallet.ID})

	sql, args, err := updateQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = p.Db.Exec(sql, args...)
	if err != nil {
		return err
	}

	selectQuery := psql.Select("id", "user_id", "user_name", "wallet_name", "wallet_type", "balance", "created_at").
		From("user_wallet").
		Where(sq.Eq{"id": wallet.ID})

	sql, args, err = selectQuery.ToSql()
	if err != nil {
		return err
	}

	err = p.Db.QueryRow(sql, args...).Scan(&wallet.ID, &wallet.UserID, &wallet.UserName, &wallet.WalletName, &wallet.WalletType, &wallet.Balance, &wallet.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteWallet(userID int) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	deleteQuery := psql.Delete("user_wallet").Where(sq.Eq{"user_id": userID})

	sql, args, err := deleteQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = p.Db.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}
