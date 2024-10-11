package util

import (
	"math/rand/v2"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type TransferRequest struct {
	ToAccount int `json:"to_account"`
	Amount    int `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	Username          string    `json:"username"`
	Number            int64     `json:"number"`
	EncryptedPassword string    `json"encypted_password"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

func (a *Account) ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
}

func NewAccount(FirstName, LastName, Username, Email, Password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		FirstName:         FirstName,
		LastName:          LastName,
		Email:             Email,
		Username:          Username,
		Number:            int64(rand.IntN(100000)),
		EncryptedPassword: string(encpw),
		Balance:           0,
		CreatedAt:         time.Now().UTC(),
	}, nil
}
