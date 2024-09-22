package testhelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/higordasneves/e-corp/pkg/domain"
)

func AssertDomainError(t *testing.T, want, got error) {
	t.Helper()

	assert.Equal(t, domain.GetErrorType(want), domain.GetErrorType(got))
	assert.Equal(t, want.Error(), got.Error())
}
