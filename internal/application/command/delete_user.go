package command

import (
	"context"

	"github.com/velonyapp/conversation/internal/domain/repo"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type DeleteUser struct {
	UserID string
}

type DeleteUserResult struct {
}

func (DeleteUser) resultType() DeleteUserResult {
	return DeleteUserResult{}
}

type DeleteUserHandler Handler[DeleteUser, DeleteUserResult]

type deleteUserHandler struct {
	userRepo repo.User
}

func NewDeleteUserHandler(
	userRepo repo.User,
) DeleteUserHandler {
	return &deleteUserHandler{
		userRepo: userRepo,
	}
}

func (h *deleteUserHandler) Handle(
	ctx context.Context,
	cmd DeleteUser,
) (DeleteUserResult, error) {
	userID := vo.NewUserID(cmd.UserID)

	if err := h.userRepo.Delete(ctx, userID); err != nil {
		return DeleteUserResult{}, err
	}

	return DeleteUserResult{}, nil
}
