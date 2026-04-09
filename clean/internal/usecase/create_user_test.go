// This file is a reference for use case layer tests.
// When creating tests for a new use case, follow this structure:
// - mock repo/infra via interface
// - test happy path, repo error, and validation error
package usecase_test

import (
	"context"
	"errors"
	"go-template/internal/contracts/infra"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
	"go-template/internal/usecase"
	"go-template/pkg/optional"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- mock ---

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(ctx context.Context, user domain.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserRepo) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *mockUserRepo) Update(ctx context.Context, user domain.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *mockUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

var _ infra.UserRepository = (*mockUserRepo)(nil)

// --- tests ---

func TestCreateUser_HappyPath(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("Create", mock.Anything, mock.AnythingOfType("domain.User")).Return(nil)

	createUC := usecase.NewCreateUser(repo)
	user, err := createUC.Execute(context.Background(), uc.CreateUserInput{
		Name:       "John",
		Surname:    optional.Some("Doe"),
		TelegramID: 123456,
	})

	require.NoError(t, err)
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, int64(123456), user.TelegramID)
	repo.AssertExpectations(t)
}

func TestCreateUser_RepoError(t *testing.T) {
	repo := new(mockUserRepo)
	repoErr := errors.New("connection refused")
	repo.On("Create", mock.Anything, mock.AnythingOfType("domain.User")).Return(repoErr)

	createUC := usecase.NewCreateUser(repo)
	_, err := createUC.Execute(context.Background(), uc.CreateUserInput{
		Name:       "John",
		TelegramID: 123456,
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}

func TestCreateUser_ValidationError(t *testing.T) {
	repo := new(mockUserRepo)

	createUC := usecase.NewCreateUser(repo)
	_, err := createUC.Execute(context.Background(), uc.CreateUserInput{
		Name:       "",
		TelegramID: 0,
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
	repo.AssertNotCalled(t, "Create")
}
