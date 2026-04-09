package delete_user

import (
	"go-template/internal/api/rest/middleware"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func New(useCase uc.DeleteUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			middleware.SetError(c, &domain.ValidationError{Fields: map[string]string{"id": "invalid format"}})
			return
		}

		if err := useCase.Execute(c.Request.Context(), id); err != nil {
			middleware.SetError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
