package rest

import (
	createnotification "go-template/internal/api/rest/create_notification"
	createuser "go-template/internal/api/rest/create_user"
	deletenotification "go-template/internal/api/rest/delete_notification"
	deleteuser "go-template/internal/api/rest/delete_user"
	getnotification "go-template/internal/api/rest/get_notification"
	getuser "go-template/internal/api/rest/get_user"
	listnotifications "go-template/internal/api/rest/list_notifications"
	listusers "go-template/internal/api/rest/list_users"
	"go-template/internal/api/rest/middleware"
	updateuser "go-template/internal/api/rest/update_user"
	uc "go-template/internal/contracts/usecase"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

type UseCases struct {
	CreateNotification uc.CreateNotification
	GetNotification    uc.GetNotification
	ListNotifications  uc.ListNotifications
	DeleteNotification uc.DeleteNotification
	CreateUser         uc.CreateUser
	GetUser            uc.GetUser
	ListUsers          uc.ListUsers
	UpdateUser         uc.UpdateUser
	DeleteUser         uc.DeleteUser
}

func NewRouter(ucs UseCases, log *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(otelgin.Middleware("rest"), gin.Logger(), gin.Recovery(), middleware.ErrorHandler(log))

	notifications := r.Group("/notifications")
	{
		notifications.POST("", createnotification.New(ucs.CreateNotification))
		notifications.GET("/:id", getnotification.New(ucs.GetNotification))
		notifications.GET("", listnotifications.New(ucs.ListNotifications))
		notifications.DELETE("/:id", deletenotification.New(ucs.DeleteNotification))
	}

	users := r.Group("/users")
	{
		users.POST("", createuser.New(ucs.CreateUser))
		users.GET("/:id", getuser.New(ucs.GetUser))
		users.GET("", listusers.New(ucs.ListUsers))
		users.PUT("/:id", updateuser.New(ucs.UpdateUser))
		users.DELETE("/:id", deleteuser.New(ucs.DeleteUser))
	}

	return r
}
