package update_user

import (
	"go-template/internal/api/rest/middleware"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
	"go-template/pkg/optional"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func New(useCase uc.UpdateUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			middleware.SetError(c, &domain.ValidationError{Fields: map[string]string{"id": "invalid format"}})
			return
		}

		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.SetError(c, &domain.ValidationError{Fields: map[string]string{"body": "invalid request format"}})
			return
		}

		input := uc.UpdateUserInput{
			ID:         id,
			Name:       req.Name,
			TelegramID: req.TelegramID,
		}
		if req.Surname != nil {
			input.Surname = optional.Some(*req.Surname)
		}

		user, err := useCase.Execute(c.Request.Context(), input)
		if err != nil {
			middleware.SetError(c, err)
			return
		}

		resp := Response{
			ID:         user.ID,
			Name:       user.Name,
			TelegramID: user.TelegramID,
		}
		if v, ok := user.Surname.Get(); ok {
			resp.Surname = &v
		}

		c.JSON(http.StatusOK, resp)
	}
}
