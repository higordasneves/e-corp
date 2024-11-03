package controller

import (
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	"github.com/higordasneves/e-corp/pkg/gateway/rabbitmq"
)

type API struct {
	AuthController
	AccountController
	TransferController
}

func NewApi(r postgres.Repository, broker rabbitmq.Publisher, cfg config.Config) API {
	accUseCase := usecase.NewAccountUseCase(r, broker)
	accController := NewAccountController(accUseCase)

	tUseCase := usecase.NewTransferUseCase(r)
	tController := NewTransferController(tUseCase)

	authUseCase := usecase.NewAuthUseCase(r, &cfg.Auth)
	authController := NewAuthController(authUseCase, cfg.Auth.SecretKey)

	return API{
		AuthController:     authController,
		AccountController:  accController,
		TransferController: tController,
	}
}
