package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
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
	return &PostgresStore{db: db}, nil
}

func (pg *PostgresStore) Init() error {
	return pg.CreateAccountTable()
}

func (pg *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)`

	_, err := pg.db.Exec(query)
	return err
}

func (pg *PostgresStore) CreateAccount(a *Account) error {
	query := `insert into account 
	(first_name, last_name, number, balance, created_at)
	values ($1, $2, $3, $4, $5)`

	resp, err := pg.db.Exec(
		query,
		a.FirstName,
		a.LastName,
		a.Number,
		a.Balance,
		a.CreatedAt,
	)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)
	return nil
}
func (pg *PostgresStore) DeleteAccount(id int) error {
	_, err := pg.db.Query("delete from account where id = $1", id)
	return err
}

func (pg *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (pg *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := pg.db.Query("select * from account where id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account id %d not found", id)
}

func (pg *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := pg.db.Query("select * from account")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		acc, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	acc := new(Account)
	err := rows.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Number,
		&acc.Balance,
		&acc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return acc, err
}
