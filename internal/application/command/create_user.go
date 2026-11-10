package command

import (
	"context"

	"github.com/velonyapp/conversation/internal/domain/entity"
	"github.com/velonyapp/conversation/internal/domain/repo"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type CreateUser struct {
	UserID string
}

type CreateUserResult struct {
}

func (CreateUser) resultType() CreateUserResult {
	return CreateUserResult{}
}

type CreateUserHandler Handler[CreateUser, CreateUserResult]

type createUserHandler struct {
	userRepo repo.User
}

func NewCreateUserHandler(
	userRepo repo.User,
) CreateUserHandler {
	return &createUserHandler{
		userRepo: userRepo,
	}
}

func (h *createUserHandler) Handle(
	ctx context.Context,
	cmd CreateUser,
) (CreateUserResult, error) {
	userID := vo.NewUserID(cmd.UserID)

	user := entity.NewUser(userID)

	if err := h.userRepo.Save(ctx, user); err != nil {
		return CreateUserResult{}, err
	}

	return CreateUserResult{}, nil
}
