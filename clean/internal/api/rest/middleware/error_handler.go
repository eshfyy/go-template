package middleware

import (
	"errors"
	"go-template/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const errKey = "handler_error"

// SetError stores an error in the gin context for ErrorHandler to process.
func SetError(c *gin.Context, err error) {
	c.Set(errKey, err)
	c.Abort()
}

// ErrorHandler is a gin middleware that maps domain errors to HTTP status codes.
// Handlers call middleware.SetError(c, err) instead of c.JSON(500, ...).
func ErrorHandler(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		val, exists := c.Get(errKey)
		if !exists {
			return
		}

		err, ok := val.(error)
		if !ok {
			return
		}

		var validationErr *domain.ValidationError
		switch {
		case errors.As(err, &validationErr):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":  "validation failed",
				"fields": validationErr.Fields,
			})
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		case errors.Is(err, domain.ErrAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "already exists"})
		case errors.Is(err, domain.ErrInvalidInput):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid input"})
		default:
			log.Error("internal error", zap.Error(err), zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
	}
}
