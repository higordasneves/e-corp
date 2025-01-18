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
	createAccUseCase := usecase.NewCreateAccountUC(r, broker)
	getAccUseCase := usecase.NewGetAccountBalanceUC(r)
	listAccUseCase := usecase.NewListAccountsUC(r)
	accountsUCs := struct {
		usecase.CreateAccountUC
		usecase.GetAccountBalanceUC
		usecase.ListAccountsUC
	}{
		createAccUseCase,
		getAccUseCase,
		listAccUseCase,
	}
	accController := NewAccountController(accountsUCs)

	tUseCase := usecase.NewTransferUC(r)
	listTransfersUC := usecase.NewListAccountTransfersUC(r)
	transfersUCs := struct {
		usecase.TransferUC
		usecase.ListAccountTransfersUC
	}{
		tUseCase,
		listTransfersUC,
	}
	tController := NewTransferController(transfersUCs)

	authUseCase := usecase.NewAuthUC(r, &cfg.Auth)
	authController := NewAuthController(authUseCase, cfg.Auth.SecretKey)

	return API{
		AuthController:     authController,
		AccountController:  accController,
		TransferController: tController,
	}
}
