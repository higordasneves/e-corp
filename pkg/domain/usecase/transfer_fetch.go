package usecase

import (
	"context"
	"fmt"
)

func (tUseCase transferUseCase) GetTransfers(ctx context.Context, cpf string) {
	fmt.Println("id:", ctx.Value("subject"))
}
