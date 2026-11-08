package service

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/command"
)

type Consumer struct {
	commandBus *command.Bus
}

func NewConsumer(
	commandBus *command.Bus,
) *Consumer {
	return &Consumer{
		commandBus: commandBus,
	}
}

type CreateUserRequest struct {
	UserID string
}

func (s *Consumer) CreateUser(ctx context.Context, req CreateUserRequest) error {
	if _, err := command.Send(ctx, s.commandBus, command.CreateUser{
		UserID: req.UserID,
	}); err != nil {
		return err
	}

	return nil
}

type DeleteUserRequest struct {
	UserID string
}

func (s *Consumer) DeleteUser(ctx context.Context, req DeleteUserRequest) error {
	if _, err := command.Send(ctx, s.commandBus, command.DeleteUser{
		UserID: req.UserID,
	}); err != nil {
		return err
	}

	return nil
}
