package internal

import (
	"database/sql"
	"fmt"

	"github.com/Nathene/MiniBank/pkg/util"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*util.Account) error
	DeleteAccount(int) error
	UpdateAccount(*util.Account) error
	GetAccounts() ([]*util.Account, error)
	GetAccountByID(int) (*util.Account, error)
	GetAccountByUsername(string) (*util.Account, error)
}

type PostgresStore struct {
	Db *sql.DB
}

func NewPostGresStore() (*PostgresStore, error) {
	// TODO: will integrate with vault soon for proper secrets managment.
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{Db: db}, nil
}

func (pg *PostgresStore) Init() error {
	return pg.CreateAccountTable()
}

func (pg *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		email varchar(50),
		number serial,
		encrypted_password varchar(100),
		balance serial,
		created_at timestamp
	)`

	_, err := pg.Db.Exec(query)
	return err
}

func (pg *PostgresStore) CreateAccount(a *util.Account) error {
	query := `insert into account 
	(first_name, last_name, email, number, encrypted_password, balance, created_at)
	values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := pg.Db.Exec(
		query,
		a.FirstName,
		a.LastName,
		a.Email,
		a.Number,
		a.EncryptedPassword,
		a.Balance,
		a.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
func (pg *PostgresStore) DeleteAccount(id int) error {
	_, err := pg.Db.Query("delete from account where id = $1", id)
	return err
}

func (pg *PostgresStore) UpdateAccount(*util.Account) error {
	return nil
}

func (pg *PostgresStore) GetAccountByUsername(username string) (*util.Account, error) {
	rows, err := pg.Db.Query("select * from account where username = $1", username)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account with number [%s] not found", username)
}

func (pg *PostgresStore) GetAccountByID(id int) (*util.Account, error) {
	rows, err := pg.Db.Query("select * from account where id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account id %d not found", id)
}

func (pg *PostgresStore) GetAccounts() ([]*util.Account, error) {
	rows, err := pg.Db.Query("select * from account")
	if err != nil {
		return nil, err
	}

	accounts := []*util.Account{}
	for rows.Next() {
		acc, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*util.Account, error) {
	acc := new(util.Account)
	err := rows.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Email,
		&acc.Number,
		&acc.EncryptedPassword,
		&acc.Balance,
		&acc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return acc, err
}
