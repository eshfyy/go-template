package get_notification

import (
	"go-template/internal/api/rest/middleware"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func New(useCase uc.GetNotification) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			middleware.SetError(c, &domain.ValidationError{Fields: map[string]string{"id": "invalid format"}})
			return
		}

		notification, err := useCase.Execute(c.Request.Context(), id)
		if err != nil {
			middleware.SetError(c, err)
			return
		}

		resp := Response{
			ID:        notification.ID,
			UserID:    notification.UserID,
			Title:     notification.Title,
			Text:      notification.Text,
			Channel:   string(notification.Channel),
			Status:    string(notification.Status),
			CreatedAt: notification.CreatedAt,
		}
		if sentAt, ok := notification.SentAt.Get(); ok {
			resp.SentAt = &sentAt
		}

		c.JSON(http.StatusOK, resp)
	}
}
