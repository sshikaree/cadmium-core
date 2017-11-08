// ------------------------------------
// TODO
// 1. Get rid of password_string !!!!!!
//
// -----------------------------------

package database

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteProvider struct {
	db *sqlx.DB
}

func NewSqliteProvider(db_filename string) (DatabaseInterface, error) {
	var err error
	provider := new(sqliteProvider)
	if provider.db, err = sqlx.Open("sqlite3", db_filename); err != nil {
		return nil, err
	}
	_, err = provider.db.Exec(
		`CREATE TABLE IF NOT EXISTS accounts
			(
				id INTEGER NOT NULL PRIMARY KEY, 
				protocol TEXT NOT NULL,
				user_name TEXT NOT NULL,
				alias TEXT DEFAULT '',
				password_string TEXT NOT NULL,
				password_hash TEXT NOT NULL,
				domain TEXT NOT NULL,
				resource TEXT DEFAULT '',
				must_encrypt INTEGER DEFAULT 0,
				port INTEGER,
				is_active INTEGER DEFAULT 1
			)`,
	)

	return provider, err
}

// Accounts section

func (p *sqliteProvider) CreateAccount(a *Account) error {
	// check for existing username@domain
	var exists bool
	err := p.db.QueryRow(
		`SELECT exists (SELECT id FROM accounts WHERE user_name = ? AND domain = ?)`,
		a.UserName, a.Domain,
	).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if exists {
		return fmt.Errorf("%s@%s already exists", a.UserName, a.Domain)
	}

	passwd_hash, err := bcrypt.GenerateFromPassword([]byte(a.PasswordString), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = p.db.Exec(
		`INSERT INTO accounts (protocol, user_name, alias, password_string, password_hash, domain, resource, must_encrypt, port)
		 VALUES (?,?,?,?,?,?,?,?,?)`,
		a.Protocol, a.UserName, a.Alias, a.PasswordString, passwd_hash, a.Domain, a.Resource, a.MustEncrypt, a.Port,
	)
	return err
}
func (p *sqliteProvider) UpdateAccount(a *Account) error {
	return nil
}
func (p *sqliteProvider) GetAccountById(id int64) (Account, error) {
	return Account{}, nil
}
func (p *sqliteProvider) GetActiveAccounts() ([]Account, error) {
	// rows, err := p.db.Query(
	// 	`SELECT id, protocol, user_name, password_hash, domain, resource, must_encrypt, port, is_active FROM accounts WHERE is_active = 1`,
	// )
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	// var accounts []Account
	// for rows.Next() {
	// 	var a *Account
	// 	err := rows.Scan(
	// 		a.ID, a.Protocol, a.UserName, a.PasswordHash, a.Domain, a.Resource, a.MustEncrypt, a.Port, a.IsActive,
	// 	)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	accounts = append(accounts, *a)

	// }

	// rewrite using sqlx
	var accounts []Account

	err := p.db.Select(
		&accounts,
		`SELECT id, protocol, user_name, password_string, password_hash, domain, resource, must_encrypt, port, is_active FROM accounts WHERE is_active = 1`,
	)
	return accounts, err
}

func (p *sqliteProvider) GetAccountsByProtocol(protocol string) ([]Account, error) {
	var accounts []Account

	err := p.db.Select(
		&accounts,
		`SELECT id, protocol, user_name, password_hash, domain, resource, must_encrypt, port, is_active FROM accounts WHERE protocol = ?`,
		protocol,
	)
	return accounts, err
}

// Contacts section

func (p *sqliteProvider) CreateContact(Contact) error {
	return nil
}
func (p *sqliteProvider) UpdateContact(Contact) error {
	return nil
}
func (p *sqliteProvider) GetContactById(id int64) (Contact, error) {
	return Contact{}, nil
}
func (p *sqliteProvider) GetContactListByAccountId(account_id int64) ([]Contact, error) {
	return []Contact{}, nil
}
