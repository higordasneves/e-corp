package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/fx"
)

func TestApp(t *testing.T) {
	t.Parallel()

	err := fx.ValidateApp(Options)
	assert.NoError(t, err)
}
