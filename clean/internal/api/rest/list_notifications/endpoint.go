package list_notifications

import (
	"go-template/internal/api/rest/middleware"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func New(useCase uc.ListNotifications) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.Query("user_id"))
		if err != nil {
			middleware.SetError(c, &domain.ValidationError{Fields: map[string]string{"user_id": "invalid format"}})
			return
		}

		limit, offset, err := parsePagination(c)
		if err != nil {
			middleware.SetError(c, err)
			return
		}

		notifications, err := useCase.Execute(c.Request.Context(), uc.ListNotificationsInput{
			UserID: userID,
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			middleware.SetError(c, err)
			return
		}

		items := make([]NotificationItem, 0, len(notifications))
		for _, n := range notifications {
			item := NotificationItem{
				ID:        n.ID,
				Title:     n.Title,
				Channel:   string(n.Channel),
				Status:    string(n.Status),
				CreatedAt: n.CreatedAt,
			}
			if sentAt, ok := n.SentAt.Get(); ok {
				item.SentAt = &sentAt
			}
			items = append(items, item)
		}

		c.JSON(http.StatusOK, Response{Items: items})
	}
}

func parsePagination(c *gin.Context) (limit, offset int, err error) {
	limit = 20
	if s := c.Query("limit"); s != "" {
		limit, err = strconv.Atoi(s)
		if err != nil {
			return 0, 0, &domain.ValidationError{Fields: map[string]string{"limit": "must be a number"}}
		}
	}
	if s := c.Query("offset"); s != "" {
		offset, err = strconv.Atoi(s)
		if err != nil {
			return 0, 0, &domain.ValidationError{Fields: map[string]string{"offset": "must be a number"}}
		}
	}
	if limit <= 0 {
		limit = 20
	}
	return limit, offset, nil
}
