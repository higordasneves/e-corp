package postgres

import (
	"database/sql"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/repository"
)

type account struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) repository.AccountRepo {
	return &account{db}
}

//CreateAccount inserts a account in database
func (a account) CreateAccount(acc *models.Account) error {
	statement, err := a.db.Prepare("INSERT INTO accounts (id, cpf, name, secret, balance, created_at) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(acc.ID, acc.CPF, acc.Name, acc.Secret, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}
