package database

type Account struct {
	ID             int64 `db:"id"`
	Protocol       string
	UserName       string `db:"user_name"`
	Alias          string
	PasswordString string `db:"password_string"`
	PasswordHash   string `db:"password_hash"`
	Domain         string
	Resource       string
	MustEncrypt    bool `db:"must_encrypt"`
	Port           int
	IsActive       bool `db:"is_active"`
}

// type XmppAccount struct {
// 	ID             int64
// 	Username       string
// 	Alias          string
// 	PasswordString string
// 	PasswordHash   string
// 	Domain         string
// 	Resource       string
// 	MustEncrypt    bool
// 	Port           int
// 	IsActive       bool
// }

type Contact struct {
	ID     int64
	Remote string
	Name   string
	Group  []string
}
type Chat struct {
	ID int64
}

type ContactStore interface {
	CreateContact(Contact) error
	UpdateContact(Contact) error
	GetContactById(contact_id int64) (Contact, error)
	GetContactListByAccountId(account_id int64) ([]Contact, error)
}
type ChatStore interface {
}

type AccountStore interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	GetAccountById(account_id int64) (Account, error)
	GetActiveAccounts() ([]Account, error)
	GetAccountsByProtocol(string) ([]Account, error)
}

type DatabaseInterface interface {
	AccountStore
	ContactStore
	// ChatStore
}

var (
	DB DatabaseInterface
)
