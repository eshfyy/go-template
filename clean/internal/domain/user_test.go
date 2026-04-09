// This file is a reference for domain layer tests.
// When creating tests for a new entity, follow this structure.
package domain_test

import (
	"errors"
	"go-template/internal/domain"
	"go-template/pkg/optional"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser_Valid(t *testing.T) {
	user, err := domain.NewUser("John", optional.Some("Doe"), 123456)

	require.NoError(t, err)
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, int64(123456), user.TelegramID)
	assert.NotEmpty(t, user.ID)
	assert.NotZero(t, user.CreatedAt)

	surname, ok := user.Surname.Get()
	assert.True(t, ok)
	assert.Equal(t, "Doe", surname)
}

func TestNewUser_WithoutSurname(t *testing.T) {
	user, err := domain.NewUser("John", optional.None[string](), 123456)

	require.NoError(t, err)
	_, ok := user.Surname.Get()
	assert.False(t, ok)
}

func TestNewUser_EmptyName(t *testing.T) {
	_, err := domain.NewUser("", optional.None[string](), 123456)

	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))

	var validationErr *domain.ValidationError
	require.ErrorAs(t, err, &validationErr)
	assert.Contains(t, validationErr.Fields, "name")
}

func TestNewUser_InvalidTelegramID(t *testing.T) {
	_, err := domain.NewUser("John", optional.None[string](), 0)

	require.Error(t, err)
	var validationErr *domain.ValidationError
	require.ErrorAs(t, err, &validationErr)
	assert.Contains(t, validationErr.Fields, "telegram_id")
}

func TestNewUser_MultipleErrors(t *testing.T) {
	_, err := domain.NewUser("", optional.None[string](), -1)

	var validationErr *domain.ValidationError
	require.ErrorAs(t, err, &validationErr)
	assert.Len(t, validationErr.Fields, 2)
	assert.Contains(t, validationErr.Fields, "name")
	assert.Contains(t, validationErr.Fields, "telegram_id")
}

func TestUser_UpdateProfile_Valid(t *testing.T) {
	user, err := domain.NewUser("John", optional.None[string](), 123)
	require.NoError(t, err)

	err = user.UpdateProfile("Jane", optional.Some("Smith"), 456)

	require.NoError(t, err)
	assert.Equal(t, "Jane", user.Name)
	assert.Equal(t, int64(456), user.TelegramID)
	surname, ok := user.Surname.Get()
	assert.True(t, ok)
	assert.Equal(t, "Smith", surname)
}

func TestUser_UpdateProfile_Invalid(t *testing.T) {
	user, err := domain.NewUser("John", optional.None[string](), 123)
	require.NoError(t, err)

	err = user.UpdateProfile("", optional.None[string](), 456)

	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Equal(t, "John", user.Name) // unchanged
}
