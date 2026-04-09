package list_users

import (
	"go-template/internal/api/rest/middleware"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func New(useCase uc.ListUsers) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, offset, err := parsePagination(c)
		if err != nil {
			middleware.SetError(c, err)
			return
		}

		users, err := useCase.Execute(c.Request.Context(), uc.ListUsersInput{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			middleware.SetError(c, err)
			return
		}

		items := make([]UserItem, len(users))
		for i, u := range users {
			items[i] = UserItem{
				ID:         u.ID,
				Name:       u.Name,
				TelegramID: u.TelegramID,
				CreatedAt:  u.CreatedAt,
			}
			if v, ok := u.Surname.Get(); ok {
				items[i].Surname = &v
			}
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
