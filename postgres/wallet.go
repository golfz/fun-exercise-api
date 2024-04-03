package postgres

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/golfz/fun-exercise-api/wallet"
	"log"
)

// type Wallet struct {
//	ID         int       `postgres:"id"`
//	UserID     int       `postgres:"user_id"`
//	UserName   string    `postgres:"user_name"`
//	WalletName string    `postgres:"wallet_name"`
//	WalletType string    `postgres:"wallet_type"`
//	Balance    float64   `postgres:"balance"`
//	CreatedAt  time.Time `postgres:"created_at"`
//}

func scanWalletFromRow(row *sql.Row) (wallet.Wallet, error) {
	var w wallet.Wallet
	err := row.Scan(&w.ID, &w.UserID, &w.UserName, &w.WalletName, &w.WalletType, &w.Balance, &w.CreatedAt)
	if err != nil {
		return wallet.Wallet{}, err
	}
	return w, nil
}

func scanWalletsFromRows(rows *sql.Rows) ([]wallet.Wallet, error) {
	wallets := make([]wallet.Wallet, 0)
	for rows.Next() {
		var w wallet.Wallet
		err := rows.Scan(&w.ID, &w.UserID, &w.UserName, &w.WalletName, &w.WalletType, &w.Balance, &w.CreatedAt)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)

	}
	return wallets, nil
}

func (p *Postgres) getWalletByID(id int) (wallet.Wallet, error) {
	selectSql := `
		SELECT id, user_id, user_name, wallet_name, wallet_type, balance, created_at 
		FROM user_wallet 
		WHERE id = $1`
	row := p.Db.QueryRow(selectSql, id)

	return scanWalletFromRow(row)
}

func prepareSelectSqlWithFilter(filter wallet.Wallet) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	selectQuery := psql.Select("id, user_id, user_name, wallet_name, wallet_type, balance, created_at").
		From("user_wallet")

	// prepare filter
	if filter.WalletType != "" {
		selectQuery = selectQuery.Where(sq.Eq{"wallet_type": filter.WalletType})
	}
	if filter.UserID != 0 {
		selectQuery = selectQuery.Where(sq.Eq{"user_id": filter.UserID})
	}

	selectQuery = selectQuery.OrderBy("id ASC")

	return selectQuery.ToSql()
}

func (p *Postgres) GetWallets(filter wallet.Wallet) ([]wallet.Wallet, error) {
	selectSql, args, err := prepareSelectSqlWithFilter(filter)
	log.Println(selectSql)
	if err != nil {
		return nil, err
	}

	rows, err := p.Db.Query(selectSql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanWalletsFromRows(rows)
}

func (p *Postgres) CreateWallet(wallet *wallet.Wallet) error {
	insertSql := `
		INSERT INTO user_wallet (user_id, user_name, wallet_name, wallet_type, balance)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	args := []interface{}{wallet.UserID, wallet.UserName, wallet.WalletName, wallet.WalletType, wallet.Balance}

	err := p.Db.QueryRow(insertSql, args...).Scan(&wallet.ID)
	if err != nil {
		return err
	}

	*wallet, err = p.getWalletByID(wallet.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateWallet(wallet *wallet.Wallet) error {
	updateSql := `UPDATE user_wallet SET balance = $1 WHERE id = $2`

	_, err := p.Db.Exec(updateSql, wallet.Balance, wallet.ID)
	if err != nil {
		return err
	}

	*wallet, err = p.getWalletByID(wallet.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteWallet(userID int) error {
	deleteSql := `DELETE FROM user_wallet WHERE user_id = $1`

	_, err := p.Db.Exec(deleteSql, userID)
	if err != nil {
		return err
	}

	return nil
}
