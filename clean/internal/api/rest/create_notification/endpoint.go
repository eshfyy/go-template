package create_notification

import (
	"go-template/internal/api/rest/middleware"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(useCase uc.CreateNotification) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.SetError(c, &domain.ValidationError{Fields: map[string]string{"body": "invalid request format"}})
			return
		}

		notification, err := useCase.Execute(c.Request.Context(), uc.CreateNotificationInput{
			UserID:  req.UserID,
			Title:   req.Title,
			Text:    req.Text,
			Channel: domain.NotificationChannel(req.Channel),
		})
		if err != nil {
			middleware.SetError(c, err)
			return
		}

		c.JSON(http.StatusCreated, Response{
			ID:      notification.ID,
			Status:  string(notification.Status),
			Channel: string(notification.Channel),
		})
	}
}
