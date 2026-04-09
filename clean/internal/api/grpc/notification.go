package grpc

import (
	"context"

	"go-template/gen/go/notificationpb"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NotificationServer struct {
	notificationpb.UnimplementedNotificationServiceServer
	createUC uc.CreateNotification
	getUC    uc.GetNotification
	listUC   uc.ListNotifications
	deleteUC uc.DeleteNotification
}

func NewNotificationServer(
	createUC uc.CreateNotification,
	getUC uc.GetNotification,
	listUC uc.ListNotifications,
	deleteUC uc.DeleteNotification,
) *NotificationServer {
	return &NotificationServer{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		deleteUC: deleteUC,
	}
}

func (s *NotificationServer) CreateNotification(ctx context.Context, req *notificationpb.CreateNotificationRequest) (*notificationpb.CreateNotificationResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	notification, err := s.createUC.Execute(ctx, uc.CreateNotificationInput{
		UserID:  userID,
		Title:   req.GetTitle(),
		Text:    req.GetText(),
		Channel: domain.NotificationChannel(req.GetChannel()),
	})
	if err != nil {
		return nil, err // interceptor maps to gRPC status
	}

	return &notificationpb.CreateNotificationResponse{
		Id:      notification.ID.String(),
		Status:  string(notification.Status),
		Channel: string(notification.Channel),
	}, nil
}

func (s *NotificationServer) GetNotification(ctx context.Context, req *notificationpb.GetNotificationRequest) (*notificationpb.GetNotificationResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	notification, err := s.getUC.Execute(ctx, id)
	if err != nil {
		return nil, err
	}

	return notificationToProto(notification), nil
}

func (s *NotificationServer) ListNotifications(ctx context.Context, req *notificationpb.ListNotificationsRequest) (*notificationpb.ListNotificationsResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	limit := int(req.GetLimit())
	if limit <= 0 {
		limit = 20
	}

	notifications, err := s.listUC.Execute(ctx, uc.ListNotificationsInput{
		UserID: userID,
		Limit:  limit,
		Offset: int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}

	items := make([]*notificationpb.GetNotificationResponse, len(notifications))
	for i, n := range notifications {
		items[i] = notificationToProto(n)
	}

	return &notificationpb.ListNotificationsResponse{Items: items}, nil
}

func (s *NotificationServer) DeleteNotification(ctx context.Context, req *notificationpb.DeleteNotificationRequest) (*notificationpb.DeleteNotificationResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	if err := s.deleteUC.Execute(ctx, id); err != nil {
		return nil, err
	}

	return &notificationpb.DeleteNotificationResponse{}, nil
}

func notificationToProto(n domain.Notification) *notificationpb.GetNotificationResponse {
	resp := &notificationpb.GetNotificationResponse{
		Id:        n.ID.String(),
		UserId:    n.UserID.String(),
		Title:     n.Title,
		Text:      n.Text,
		Channel:   string(n.Channel),
		Status:    string(n.Status),
		CreatedAt: timestamppb.New(n.CreatedAt),
	}
	if sentAt, ok := n.SentAt.Get(); ok {
		resp.SentAt = timestamppb.New(sentAt)
	}
	return resp
}
