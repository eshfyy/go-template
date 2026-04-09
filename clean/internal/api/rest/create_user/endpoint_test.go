// This file is a reference for transport layer tests.
// When creating tests for a new endpoint, follow this structure:
// - mock use case via interface
// - test happy path (201), validation error (422), and internal error (500)
package create_user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	createuser "go-template/internal/api/rest/create_user"
	"go-template/internal/api/rest/middleware"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
	"go-template/pkg/optional"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// --- mock ---

type mockCreateUser struct {
	mock.Mock
}

func (m *mockCreateUser) Execute(ctx context.Context, input uc.CreateUserInput) (domain.User, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(domain.User), args.Error(1)
}

var _ uc.CreateUser = (*mockCreateUser)(nil)

// --- helpers ---

func setupRouter(ucMock uc.CreateUser) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler(zap.NewNop()))
	r.POST("/users", createuser.New(ucMock))
	return r
}

func doRequest(t *testing.T, router *gin.Engine, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// --- tests ---

func TestCreateUser_201(t *testing.T) {
	ucMock := new(mockCreateUser)
	user, _ := domain.NewUser("John", optional.Some("Doe"), 123456)
	ucMock.On("Execute", mock.Anything, mock.AnythingOfType("usecase.CreateUserInput")).Return(user, nil)

	router := setupRouter(ucMock)
	w := doRequest(t, router, map[string]any{
		"name":        "John",
		"surname":     "Doe",
		"telegram_id": 123456,
	})

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp createuser.Response
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "John", resp.Name)
	ucMock.AssertExpectations(t)
}

func TestCreateUser_422_ValidationError(t *testing.T) {
	ucMock := new(mockCreateUser)
	ucMock.On("Execute", mock.Anything, mock.Anything).Return(
		domain.User{},
		&domain.ValidationError{Fields: map[string]string{"telegram_id": "must be positive"}},
	)

	router := setupRouter(ucMock)
	// gin binding passes, but use case returns domain validation error
	w := doRequest(t, router, map[string]any{
		"name":        "John",
		"telegram_id": -1,
	})

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Contains(t, body, "fields")
}

func TestCreateUser_500_InternalError(t *testing.T) {
	ucMock := new(mockCreateUser)
	ucMock.On("Execute", mock.Anything, mock.Anything).Return(
		domain.User{},
		errors.New("db connection lost"),
	)

	router := setupRouter(ucMock)
	w := doRequest(t, router, map[string]any{
		"name":        "John",
		"telegram_id": 123,
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "internal error", body["error"]) // no leak
}
