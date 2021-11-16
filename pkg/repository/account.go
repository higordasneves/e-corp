package repository

import (
	"github.com/higordasneves/e-corp/pkg/domain/models"
)

type AccountRepo interface {
	CreateAccount(account models.Account) error
}
