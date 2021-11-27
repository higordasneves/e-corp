package postgres

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type transfer struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewTransferRepository(dbPool *pgxpool.Pool, log *logrus.Logger) repository.TransferRepo {
	return &transfer{db: dbPool, log: log}
}

func (t transfer) Transfer(ctx context.Context, accOriID vos.UUID, accDestID vos.UUID, amount vos.Currency) {

}

func (t transfer) GetTransfers(ctx context.Context, cpf string) {

}
