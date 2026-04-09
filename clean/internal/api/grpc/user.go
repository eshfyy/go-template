package grpc

import (
	"context"

	"go-template/gen/go/userpb"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
	"go-template/pkg/optional"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	createUC uc.CreateUser
	getUC    uc.GetUser
	listUC   uc.ListUsers
	updateUC uc.UpdateUser
	deleteUC uc.DeleteUser
}

func NewUserServer(
	createUC uc.CreateUser,
	getUC uc.GetUser,
	listUC uc.ListUsers,
	updateUC uc.UpdateUser,
	deleteUC uc.DeleteUser,
) *UserServer {
	return &UserServer{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}
}

func (s *UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	input := uc.CreateUserInput{
		Name:       req.GetName(),
		TelegramID: req.GetTelegramId(),
	}
	if req.Surname != nil {
		input.Surname = optional.Some(*req.Surname)
	}

	user, err := s.createUC.Execute(ctx, input)
	if err != nil {
		return nil, err
	}

	resp := &userpb.CreateUserResponse{
		Id:         user.ID.String(),
		Name:       user.Name,
		TelegramId: user.TelegramID,
	}
	if v, ok := user.Surname.Get(); ok {
		resp.Surname = &v
	}
	return resp, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	user, err := s.getUC.Execute(ctx, id)
	if err != nil {
		return nil, err
	}

	return userToProto(user), nil
}

func (s *UserServer) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	limit := int(req.GetLimit())
	if limit <= 0 {
		limit = 20
	}

	users, err := s.listUC.Execute(ctx, uc.ListUsersInput{
		Limit:  limit,
		Offset: int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}

	items := make([]*userpb.GetUserResponse, len(users))
	for i, u := range users {
		items[i] = userToProto(u)
	}

	return &userpb.ListUsersResponse{Items: items}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	input := uc.UpdateUserInput{
		ID:         id,
		Name:       req.GetName(),
		TelegramID: req.GetTelegramId(),
	}
	if req.Surname != nil {
		input.Surname = optional.Some(*req.Surname)
	}

	user, err := s.updateUC.Execute(ctx, input)
	if err != nil {
		return nil, err
	}

	resp := &userpb.UpdateUserResponse{
		Id:         user.ID.String(),
		Name:       user.Name,
		TelegramId: user.TelegramID,
	}
	if v, ok := user.Surname.Get(); ok {
		resp.Surname = &v
	}
	return resp, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	if err := s.deleteUC.Execute(ctx, id); err != nil {
		return nil, err
	}

	return &userpb.DeleteUserResponse{}, nil
}

func userToProto(u domain.User) *userpb.GetUserResponse {
	resp := &userpb.GetUserResponse{
		Id:         u.ID.String(),
		Name:       u.Name,
		TelegramId: u.TelegramID,
		CreatedAt:  timestamppb.New(u.CreatedAt),
	}
	if v, ok := u.Surname.Get(); ok {
		resp.Surname = &v
	}
	return resp
}
