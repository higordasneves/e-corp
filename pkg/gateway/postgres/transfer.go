package postgres

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

type transfer struct {
	dbPool *pgxpool.Pool
	log    *logrus.Logger
}

func NewTransferRepository(dbPool *pgxpool.Pool, log *logrus.Logger) repository.TransferRepo {
	return &transfer{dbPool: dbPool, log: log}
}

func (t transfer) Transfer(ctx context.Context, accOriID vos.UUID, accDestID vos.UUID, amount vos.Currency) {
	transferID := vos.NewAccID()
	_, err := t.dbPool.Exec(ctx, "INSERT INTO transfers "+
		"(id, account_origin_id, account_destination_id, amount, created_at)"+
		" VALUES ($1, $2, $3, $4, $5, $6)", transferID, accOriID, accDestID, int(amount), time.Now())
	if err != nil {
		return
	}
}

func (t transfer) GetTransfers(ctx context.Context, cpf string) {

}
