package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetAccountByID(id int) (*Account, error)
	GetAccountByNumber(number int) (*Account, error)
	CreateAccount(a *Account) error
	DeleteAccount(id int) error
	UpdateAccount(a *Account) error
	GetAccounts() ([]*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}
//sql password: murphgobank
func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "user=postgres dbname=postgres password=murphgobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Init() error{
	return s.CreateAccountTable()
}

func (s *PostgresStorage) CreateAccountTable() error{
	query:=`CREATE TABLE IF NOT EXISTS account(
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		number BIGINT,
		balance BIGINT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err:=s.db.Exec(query)
	return err
}

func (s *PostgresStorage) GetAccountByID(id int) (*Account, error) {
	
	rows,err:= s.db.Query("SELECT * FROM account WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return s.ScanAccount(rows)
	}
	return nil, fmt.Errorf("Account %d not found",id)

}

func (s *PostgresStorage) CreateAccount(a *Account) error {

    fmt.Printf("Attempting to create account: %+v\n", a)
    query := `
        INSERT INTO account (
            first_name, 
            last_name, 
            number, 
            balance, 
            created_at
        ) VALUES ($1, $2, $3, $4, $5)` 
	resp, err := s.db.Query(query, a.FirstName, a.LastName, a.Number, a.Balance, a.CreatedAt)
    if err != nil {
		fmt.Printf("Error creating account: %v", err)
        return err
    }

    fmt.Printf("resp: %+v\n", resp)
    return nil
}
func (s *PostgresStorage) GetAccountByNumber(number int) (*Account,error) {
	rows, err := s.db.Query("SELECT * FROM account WHERE number = $1",number)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return s.ScanAccount(rows)
	}
	return nil, fmt.Errorf("Account %d not found",number)
}
func (s *PostgresStorage) DeleteAccount(id int) error {
	_, err := s.db.Exec("DELETE FROM account WHERE id = $1", id)
	if err != nil {
		fmt.Printf("Error deleting account: %v", err)
		return err
	}
	return nil
}

func (s *PostgresStorage) UpdateAccount(a *Account) error {
	return nil
}

func (s *PostgresStorage) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, err
	}
	accounts := make([]*Account, 0)
	for rows.Next() {
		var a Account
		if err := rows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.Number, &a.Balance, &a.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, &a)
	}
	return accounts, nil
}

func (s *PostgresStorage) ScanAccount(rows *sql.Rows) (*Account, error) {
	var a Account
	if err := rows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.Number, &a.Balance, &a.CreatedAt); err != nil {
		return nil, err
	}
	return &a, nil
}