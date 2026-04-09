package get_user

import (
	"go-template/internal/api/rest/middleware"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func New(useCase uc.GetUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			middleware.SetError(c, &domain.ValidationError{Fields: map[string]string{"id": "invalid format"}})
			return
		}

		user, err := useCase.Execute(c.Request.Context(), id)
		if err != nil {
			middleware.SetError(c, err)
			return
		}

		resp := Response{
			ID:         user.ID,
			Name:       user.Name,
			TelegramID: user.TelegramID,
			CreatedAt:  user.CreatedAt,
		}
		if v, ok := user.Surname.Get(); ok {
			resp.Surname = &v
		}

		c.JSON(http.StatusOK, resp)
	}
}
