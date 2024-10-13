// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/controller"
	"sync"
)

// Ensure, that TransferUseCaseMock does implement controller.TransferUseCase.
// If this is not the case, regenerate this file with moq.
var _ controller.TransferUseCase = &TransferUseCaseMock{}

// TransferUseCaseMock is a mock implementation of controller.TransferUseCase.
//
//	func TestSomethingThatUsesTransferUseCase(t *testing.T) {
//
//		// make and configure a mocked controller.TransferUseCase
//		mockedTransferUseCase := &TransferUseCaseMock{
//			ListAccountTransfersFunc: func(ctx context.Context, input usecase.ListAccountTransfersInput) (usecase.ListAccountTransfersOutput, error) {
//				panic("mock out the ListAccountTransfers method")
//			},
//			TransferFunc: func(ctx context.Context, input usecase.TransferInput) (usecase.TransferOutput, error) {
//				panic("mock out the Transfer method")
//			},
//		}
//
//		// use mockedTransferUseCase in code that requires controller.TransferUseCase
//		// and then make assertions.
//
//	}
type TransferUseCaseMock struct {
	// ListAccountTransfersFunc mocks the ListAccountTransfers method.
	ListAccountTransfersFunc func(ctx context.Context, input usecase.ListAccountTransfersInput) (usecase.ListAccountTransfersOutput, error)

	// TransferFunc mocks the Transfer method.
	TransferFunc func(ctx context.Context, input usecase.TransferInput) (usecase.TransferOutput, error)

	// calls tracks calls to the methods.
	calls struct {
		// ListAccountTransfers holds details about calls to the ListAccountTransfers method.
		ListAccountTransfers []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Input is the input argument value.
			Input usecase.ListAccountTransfersInput
		}
		// Transfer holds details about calls to the Transfer method.
		Transfer []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Input is the input argument value.
			Input usecase.TransferInput
		}
	}
	lockListAccountTransfers sync.RWMutex
	lockTransfer             sync.RWMutex
}

// ListAccountTransfers calls ListAccountTransfersFunc.
func (mock *TransferUseCaseMock) ListAccountTransfers(ctx context.Context, input usecase.ListAccountTransfersInput) (usecase.ListAccountTransfersOutput, error) {
	callInfo := struct {
		Ctx   context.Context
		Input usecase.ListAccountTransfersInput
	}{
		Ctx:   ctx,
		Input: input,
	}
	mock.lockListAccountTransfers.Lock()
	mock.calls.ListAccountTransfers = append(mock.calls.ListAccountTransfers, callInfo)
	mock.lockListAccountTransfers.Unlock()
	if mock.ListAccountTransfersFunc == nil {
		var (
			listAccountTransfersOutputOut usecase.ListAccountTransfersOutput
			errOut                        error
		)
		return listAccountTransfersOutputOut, errOut
	}
	return mock.ListAccountTransfersFunc(ctx, input)
}

// ListAccountTransfersCalls gets all the calls that were made to ListAccountTransfers.
// Check the length with:
//
//	len(mockedTransferUseCase.ListAccountTransfersCalls())
func (mock *TransferUseCaseMock) ListAccountTransfersCalls() []struct {
	Ctx   context.Context
	Input usecase.ListAccountTransfersInput
} {
	var calls []struct {
		Ctx   context.Context
		Input usecase.ListAccountTransfersInput
	}
	mock.lockListAccountTransfers.RLock()
	calls = mock.calls.ListAccountTransfers
	mock.lockListAccountTransfers.RUnlock()
	return calls
}

// Transfer calls TransferFunc.
func (mock *TransferUseCaseMock) Transfer(ctx context.Context, input usecase.TransferInput) (usecase.TransferOutput, error) {
	callInfo := struct {
		Ctx   context.Context
		Input usecase.TransferInput
	}{
		Ctx:   ctx,
		Input: input,
	}
	mock.lockTransfer.Lock()
	mock.calls.Transfer = append(mock.calls.Transfer, callInfo)
	mock.lockTransfer.Unlock()
	if mock.TransferFunc == nil {
		var (
			transferOutputOut usecase.TransferOutput
			errOut            error
		)
		return transferOutputOut, errOut
	}
	return mock.TransferFunc(ctx, input)
}

// TransferCalls gets all the calls that were made to Transfer.
// Check the length with:
//
//	len(mockedTransferUseCase.TransferCalls())
func (mock *TransferUseCaseMock) TransferCalls() []struct {
	Ctx   context.Context
	Input usecase.TransferInput
} {
	var calls []struct {
		Ctx   context.Context
		Input usecase.TransferInput
	}
	mock.lockTransfer.RLock()
	calls = mock.calls.Transfer
	mock.lockTransfer.RUnlock()
	return calls
}
