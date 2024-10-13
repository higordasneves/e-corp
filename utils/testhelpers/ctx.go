package testhelpers

import (
	"context"
	"testing"
)

func NewCtx(t *testing.T) context.Context {
	return context.Background()
}
