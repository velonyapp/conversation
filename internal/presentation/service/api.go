package service

import (
	"context"
	"errors"

	v1 "github.com/velonyapp/conversation/gen/api/v1"
	"github.com/velonyapp/conversation/internal/application/command"
	"github.com/velonyapp/conversation/internal/application/query"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	kjwt "github.com/go-kratos/kratos/contrib/middleware/jwt/v3"
	"go.einride.tech/aip/resourcename"
)

var (
	ErrInvalidConversationResourceName = errors.New("invalid conversation resource name")
	ErrInvalidConversationBatchSize    = errors.New("cannot fetch more than 100 conversations at once")

	ErrInvalidConversationUserResourceName = errors.New("invalid conversation user resource name")
	ErrInvalidConversationUserBatchSize    = errors.New("cannot fetch more than 100 conversation users at once")

	ErrInvalidMessageResourceName = errors.New("invalid message resource name")
	ErrInvalidMessageBatchSize    = errors.New("cannot fetch more than 100 messages at once")

	ErrInvalidUserResourceName = errors.New("invalid user resource name")

	ErrInvalidUpdateMaskPath = errors.New("invalid update mask path")
	ErrUnauthenticated       = errors.New("unauthenticated")
	ErrForbidden             = errors.New("unauthorized")
)

const (
	conversationResourcePattern  = "conversations/{conversation}"
	conversationDefaultPageSize  = 20
	conversationMaximumPageSize  = 100
	conversationMaximumBatchSize = 100

	conversationUserResourcePattern  = "conversations/{conversation}/users/{conversation_user}"
	conversationUserDefaultPageSize  = 20
	conversationUserMaximumPageSize  = 100
	conversationUserMaximumBatchSize = 100

	messageResourcePattern  = "conversations/{conversation}/messages/{message}"
	messageDefaultPageSize  = 50
	messageMaximumPageSize  = 100
	messageMaximumBatchSize = 100

	userResourcePattern = "users/{user}"
)

type API struct {
	v1.UnimplementedConversationServiceServer

	commandBus *command.Bus
	queryBus   *query.Bus
}

func NewAPI(
	commandBus *command.Bus,
	queryBus *query.Bus,
) *API {
	return &API{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

func (s *API) GetConversation(ctx context.Context, req *v1.GetConversationRequest) (*v1.Conversation, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, err := parseConversationName(req.Name)
	if err != nil {
		return nil, err
	}

	result, err := query.Send(ctx, s.queryBus, query.GetConversation{
		ActorID:        subject,
		ConversationID: conversationID,
	})
	if err != nil {
		return nil, err
	}

	c := &v1.Conversation{
		Name:        formatConversationName(result.Conversation.ID),
		Title:       result.Conversation.Title,
		Description: result.Conversation.Description,
		CreateTime:  timestamppb.New(result.Conversation.CreateTime),
		UpdateTime:  timestamppb.New(result.Conversation.UpdateTime),
	}
	if lastMessageID := result.Conversation.LastMessageID; lastMessageID != nil {
		value := formatMessageName(result.Conversation.ID, *lastMessageID)
		c.LastMessage = &value
	}

	return c, nil
}

func (s *API) ListConversations(ctx context.Context, req *v1.ListConversationsRequest) (*v1.ListConversationsResponse, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	if req.PageSize == 0 {
		req.PageSize = conversationDefaultPageSize
	}
	req.PageSize = min(req.PageSize, conversationMaximumPageSize)

	result, err := query.Send(ctx, s.queryBus, query.ListConversations{
		ActorID:   subject,
		PageSize:  req.PageSize,
		PageToken: req.PageToken,
	})
	if err != nil {
		return nil, err
	}

	conversations := make([]*v1.Conversation, len(result.Conversations))
	for i, conversation := range result.Conversations {
		conversations[i] = &v1.Conversation{
			Name:        formatConversationName(conversation.ID),
			Title:       conversation.Title,
			Description: conversation.Description,
			CreateTime:  timestamppb.New(conversation.CreateTime),
			UpdateTime:  timestamppb.New(conversation.UpdateTime),
		}
		if lastMessageID := conversation.LastMessageID; lastMessageID != nil {
			value := formatMessageName(conversation.ID, *lastMessageID)
			conversations[i].LastMessage = &value
		}
	}

	return &v1.ListConversationsResponse{
		Conversations: conversations,
		NextPageToken: result.NextPageToken,
	}, nil
}

func (s *API) BatchGetConversations(ctx context.Context, req *v1.BatchGetConversationsRequest) (*v1.BatchGetConversationsResponse, error) {
	if len(req.Names) > conversationMaximumBatchSize {
		return nil, ErrInvalidConversationBatchSize
	}

	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationIDs := make([]string, len(req.Names))
	for i, name := range req.Names {
		conversationID, err := parseConversationName(name)
		if err != nil {
			return nil, err
		}

		conversationIDs[i] = conversationID
	}

	result, err := query.Send(ctx, s.queryBus, query.BatchGetConversations{
		ActorID:         subject,
		ConversationIDs: conversationIDs,
	})
	if err != nil {
		return nil, err
	}

	conversations := make([]*v1.Conversation, len(result.Conversations))
	for i, conversation := range result.Conversations {
		conversations[i] = &v1.Conversation{
			Name:        formatConversationName(conversation.ID),
			Title:       conversation.Title,
			Description: conversation.Description,
			CreateTime:  timestamppb.New(conversation.CreateTime),
			UpdateTime:  timestamppb.New(conversation.UpdateTime),
		}
		if lastMessageID := conversation.LastMessageID; lastMessageID != nil {
			value := formatMessageName(conversation.ID, *lastMessageID)
			conversations[i].LastMessage = &value
		}
	}

	return &v1.BatchGetConversationsResponse{
		Conversations: conversations,
	}, nil
}

func (s *API) CreateConversation(ctx context.Context, req *v1.CreateConversationRequest) (*v1.Conversation, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	result, err := command.Send(ctx, s.commandBus, command.CreateConversation{
		ActorID:     subject,
		Title:       req.Conversation.Title,
		Description: req.Conversation.Description,
	})
	if err != nil {
		return nil, err
	}

	return &v1.Conversation{
		Name:        formatConversationName(result.Conversation.ID),
		Title:       result.Conversation.Title,
		Description: result.Conversation.Description,
		CreateTime:  timestamppb.New(result.Conversation.CreateTime),
		UpdateTime:  timestamppb.New(result.Conversation.UpdateTime),
	}, nil
}

func (s *API) UpdateConversation(ctx context.Context, req *v1.UpdateConversationRequest) (*v1.Conversation, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, err := parseConversationName(req.Conversation.Name)
	if err != nil {
		return nil, err
	}

	var title, description *string
	if req.UpdateMask == nil || req.UpdateMask.Paths == nil {
		if req.Conversation.Title != "" {
			value := req.Conversation.Title
			title = &value
		}
		if req.Conversation.Description != "" {
			value := req.Conversation.Description
			description = &value
		}
	} else {
		for _, path := range req.UpdateMask.Paths {
			switch path {
			case "*":
				titleValue := req.Conversation.Title
				title = &titleValue

				descriptionValue := req.Conversation.Description
				description = &descriptionValue

			case "title":
				value := req.Conversation.Title
				title = &value

			case "description":
				value := req.Conversation.Description
				description = &value

			default:
				return nil, ErrInvalidUpdateMaskPath
			}
		}
	}

	result, err := command.Send(ctx, s.commandBus, command.UpdateConversation{
		ActorID:        subject,
		ConversationID: conversationID,
		Title:          title,
		Description:    description,
	})
	if err != nil {
		return nil, err
	}

	c := &v1.Conversation{
		Name:        formatConversationName(result.Conversation.ID),
		Title:       result.Conversation.Title,
		Description: result.Conversation.Description,
		CreateTime:  timestamppb.New(result.Conversation.CreateTime),
		UpdateTime:  timestamppb.New(result.Conversation.UpdateTime),
	}
	if lastMessageID := result.Conversation.LastMessageID; lastMessageID != nil {
		value := formatMessageName(result.Conversation.ID, *lastMessageID)
		c.LastMessage = &value
	}

	return c, nil
}

func (s *API) DeleteConversation(ctx context.Context, req *v1.DeleteConversationRequest) (*emptypb.Empty, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, err := parseConversationName(req.Name)
	if err != nil {
		return nil, err
	}

	if _, err := command.Send(ctx, s.commandBus, command.DeleteConversation{
		ActorID:        subject,
		ConversationID: conversationID,
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *API) JoinConversation(ctx context.Context, req *v1.JoinConversationRequest) (*v1.Conversation, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, err := parseConversationName(req.Name)
	if err != nil {
		return nil, err
	}

	result, err := command.Send(ctx, s.commandBus, command.JoinConversation{
		ActorID:        subject,
		ConversationID: conversationID,
	})
	if err != nil {
		return nil, err
	}

	c := &v1.Conversation{
		Name:        formatConversationName(result.Conversation.ID),
		Title:       result.Conversation.Title,
		Description: result.Conversation.Description,
		CreateTime:  timestamppb.New(result.Conversation.CreateTime),
		UpdateTime:  timestamppb.New(result.Conversation.UpdateTime),
	}
	if lastMessageID := result.Conversation.LastMessageID; lastMessageID != nil {
		value := formatMessageName(result.Conversation.ID, *lastMessageID)
		c.LastMessage = &value
	}

	return c, nil
}

func (s *API) LeaveConversation(ctx context.Context, req *v1.LeaveConversationRequest) (*v1.LeaveConversationResponse, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, err := parseConversationName(req.Name)
	if err != nil {
		return nil, err
	}

	if _, err := command.Send(ctx, s.commandBus, command.LeaveConversation{
		ActorID:        subject,
		ConversationID: conversationID,
	}); err != nil {
		return nil, err
	}

	return &v1.LeaveConversationResponse{}, nil
}

func (s *API) GetConversationUser(ctx context.Context, req *v1.GetConversationUserRequest) (*v1.ConversationUser, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, conversationUserID, err := parseConversationUserName(req.Name)
	if err != nil {
		return nil, err
	}

	result, err := query.Send(ctx, s.queryBus, query.GetConversationUser{
		ActorID:            subject,
		ConversationID:     conversationID,
		ConversationUserID: conversationUserID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ConversationUser{
		Name:       formatConversationUserName(conversationID, result.ConversationUser.ID),
		User:       formatUserName(result.ConversationUser.UserID),
		Role:       formatConversationUserRoleAPIV1(result.ConversationUser.Role),
		CreateTime: timestamppb.New(result.ConversationUser.CreateTime),
		UpdateTime: timestamppb.New(result.ConversationUser.UpdateTime),
	}, nil
}

func (s *API) ListConversationUsers(ctx context.Context, req *v1.ListConversationUsersRequest) (*v1.ListConversationUsersResponse, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, err := parseConversationName(req.Parent)
	if err != nil {
		return nil, err
	}

	if req.PageSize == 0 {
		req.PageSize = conversationUserDefaultPageSize
	}
	req.PageSize = min(req.PageSize, conversationUserMaximumPageSize)

	result, err := query.Send(ctx, s.queryBus, query.ListConversationUsers{
		ActorID:        subject,
		ConversationID: conversationID,
		PageSize:       req.PageSize,
		PageToken:      req.PageToken,
	})
	if err != nil {
		return nil, err
	}

	conversationUsers := make([]*v1.ConversationUser, len(result.ConversationUsers))
	for i, conversationUser := range result.ConversationUsers {
		conversationUsers[i] = &v1.ConversationUser{
			Name:       formatConversationUserName(conversationID, conversationUser.ID),
			User:       formatUserName(conversationUser.UserID),
			Role:       formatConversationUserRoleAPIV1(conversationUser.Role),
			CreateTime: timestamppb.New(conversationUser.CreateTime),
			UpdateTime: timestamppb.New(conversationUser.UpdateTime),
		}
	}

	return &v1.ListConversationUsersResponse{
		ConversationUsers: conversationUsers,
		NextPageToken:     result.NextPageToken,
	}, nil
}

func (s *API) BatchGetConversationUsers(ctx context.Context, req *v1.BatchGetConversationUsersRequest) (*v1.BatchGetConversationUsersResponse, error) {
	if len(req.Names) > conversationUserMaximumBatchSize {
		return nil, ErrInvalidConversationUserBatchSize
	}

	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, err := parseConversationName(req.Parent)
	if err != nil {
		return nil, err
	}
	conversationUserIDs := make([]string, len(req.Names))
	for i, name := range req.Names {
		_, conversationUserID, err := parseConversationUserName(name)
		if err != nil {
			return nil, err
		}

		conversationUserIDs[i] = conversationUserID
	}

	result, err := query.Send(ctx, s.queryBus, query.BatchGetConversationUsers{
		ActorID:             subject,
		ConversationID:      conversationID,
		ConversationUserIDs: conversationUserIDs,
	})
	if err != nil {
		return nil, err
	}

	conversationUsers := make([]*v1.ConversationUser, len(result.ConversationUsers))
	for i, conversationUser := range result.ConversationUsers {
		conversationUsers[i] = &v1.ConversationUser{
			Name:       formatConversationUserName(conversationID, conversationUser.ID),
			User:       formatUserName(conversationUser.UserID),
			Role:       formatConversationUserRoleAPIV1(conversationUser.Role),
			CreateTime: timestamppb.New(conversationUser.CreateTime),
			UpdateTime: timestamppb.New(conversationUser.UpdateTime),
		}
	}

	return &v1.BatchGetConversationUsersResponse{
		ConversationUsers: conversationUsers,
	}, nil
}

func (s *API) DeleteConversationUser(ctx context.Context, req *v1.DeleteConversationUserRequest) (*emptypb.Empty, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, conversationUserID, err := parseConversationUserName(req.Name)
	if err != nil {
		return nil, err
	}

	if _, err := command.Send(ctx, s.commandBus, command.KickConversationUser{
		ActorID:        subject,
		ConversationID: conversationID,
		TargetID:       conversationUserID,
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *API) CreateConversationUser(ctx context.Context, req *v1.CreateConversationUserRequest) (*v1.ConversationUser, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, err := parseConversationName(req.Parent)
	if err != nil {
		return nil, err
	}
	userID, err := parseUserName(req.ConversationUser.User)
	if err != nil {
		return nil, err
	}

	result, err := command.Send(ctx, s.commandBus, command.AddConversationUser{
		ActorID:        subject,
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ConversationUser{
		Name:       formatConversationUserName(conversationID, result.ConversationUser.ID),
		User:       formatUserName(result.ConversationUser.UserID),
		Role:       formatConversationUserRoleAPIV1(result.ConversationUser.Role),
		CreateTime: timestamppb.New(result.ConversationUser.CreateTime),
		UpdateTime: timestamppb.New(result.ConversationUser.UpdateTime),
	}, nil
}

func (s *API) ElectConversationUser(ctx context.Context, req *v1.ElectConversationUserRequest) (*v1.ConversationUser, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, conversationUserID, err := parseConversationUserName(req.Name)
	if err != nil {
		return nil, err
	}

	result, err := command.Send(ctx, s.commandBus, command.ElectConversationUser{
		ActorID:        subject,
		ConversationID: conversationID,
		TargetID:       conversationUserID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ConversationUser{
		Name:       formatConversationUserName(conversationID, result.ConversationUser.ID),
		User:       formatUserName(result.ConversationUser.UserID),
		Role:       formatConversationUserRoleAPIV1(result.ConversationUser.Role),
		CreateTime: timestamppb.New(result.ConversationUser.CreateTime),
		UpdateTime: timestamppb.New(result.ConversationUser.UpdateTime),
	}, nil
}

func (s *API) PromoteConversationUser(ctx context.Context, req *v1.PromoteConversationUserRequest) (*v1.ConversationUser, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, conversationUserID, err := parseConversationUserName(req.Name)
	if err != nil {
		return nil, err
	}

	result, err := command.Send(ctx, s.commandBus, command.PromoteConversationUser{
		ActorID:        subject,
		ConversationID: conversationID,
		TargetID:       conversationUserID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ConversationUser{
		Name:       formatConversationUserName(conversationID, result.ConversationUser.ID),
		User:       formatUserName(result.ConversationUser.UserID),
		Role:       formatConversationUserRoleAPIV1(result.ConversationUser.Role),
		CreateTime: timestamppb.New(result.ConversationUser.CreateTime),
		UpdateTime: timestamppb.New(result.ConversationUser.UpdateTime),
	}, nil
}

func (s *API) DemoteConversationUser(ctx context.Context, req *v1.DemoteConversationUserRequest) (*v1.ConversationUser, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, conversationUserID, err := parseConversationUserName(req.Name)
	if err != nil {
		return nil, err
	}

	result, err := command.Send(ctx, s.commandBus, command.DemoteConversationUser{
		ActorID:        subject,
		ConversationID: conversationID,
		TargetID:       conversationUserID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ConversationUser{
		Name:       formatConversationUserName(conversationID, result.ConversationUser.ID),
		User:       formatUserName(result.ConversationUser.UserID),
		Role:       formatConversationUserRoleAPIV1(result.ConversationUser.Role),
		CreateTime: timestamppb.New(result.ConversationUser.CreateTime),
		UpdateTime: timestamppb.New(result.ConversationUser.UpdateTime),
	}, nil
}

func (s *API) GetMessage(ctx context.Context, req *v1.GetMessageRequest) (*v1.Message, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, messageID, err := parseMessageName(req.Name)
	if err != nil {
		return nil, err
	}

	result, err := query.Send(ctx, s.queryBus, query.GetMessage{
		ActorID:        subject,
		ConversationID: conversationID,
		MessageID:      messageID,
	})
	if err != nil {
		return nil, err
	}

	m := &v1.Message{
		Name:       formatMessageName(conversationID, result.Message.ID),
		Content:    result.Message.Content,
		CreateTime: timestamppb.New(result.Message.CreateTime),
		UpdateTime: timestamppb.New(result.Message.UpdateTime),
	}
	if userID := result.Message.UserID; userID != nil {
		value := formatUserName(*userID)
		m.User = &value
	}

	return m, nil
}

func (s *API) ListMessages(ctx context.Context, req *v1.ListMessagesRequest) (*v1.ListMessagesResponse, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, err := parseConversationName(req.Parent)
	if err != nil {
		return nil, err
	}

	if req.PageSize == 0 {
		req.PageSize = messageDefaultPageSize
	}
	req.PageSize = min(req.PageSize, messageMaximumPageSize)

	result, err := query.Send(ctx, s.queryBus, query.ListMessages{
		ActorID:        subject,
		ConversationID: conversationID,
		PageSize:       req.PageSize,
		PageToken:      req.PageToken,
	})
	if err != nil {
		return nil, err
	}

	messages := make([]*v1.Message, len(result.Messages))
	for i, message := range result.Messages {
		messages[i] = &v1.Message{
			Name:       formatMessageName(conversationID, message.ID),
			Content:    message.Content,
			CreateTime: timestamppb.New(message.CreateTime),
			UpdateTime: timestamppb.New(message.UpdateTime),
		}
		if userID := message.UserID; userID != nil {
			value := formatUserName(*userID)
			messages[i].User = &value
		}
	}

	return &v1.ListMessagesResponse{
		Messages:      messages,
		NextPageToken: result.NextPageToken,
	}, nil
}

func (s *API) BatchGetMessages(ctx context.Context, req *v1.BatchGetMessagesRequest) (*v1.BatchGetMessagesResponse, error) {
	if len(req.Names) > messageMaximumBatchSize {
		return nil, ErrInvalidMessageBatchSize
	}

	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, err := parseConversationName(req.Parent)
	if err != nil {
		return nil, err
	}
	messageIDs := make([]string, len(req.Names))
	for i, name := range req.Names {
		_, messageID, err := parseMessageName(name)
		if err != nil {
			return nil, err
		}

		messageIDs[i] = messageID
	}

	result, err := query.Send(ctx, s.queryBus, query.BatchGetMessages{
		ActorID:        subject,
		ConversationID: conversationID,
		MessageIDs:     messageIDs,
	})
	if err != nil {
		return nil, err
	}

	messages := make([]*v1.Message, len(result.Messages))
	for i, message := range result.Messages {
		messages[i] = &v1.Message{
			Name:       formatMessageName(conversationID, message.ID),
			Content:    message.Content,
			CreateTime: timestamppb.New(message.CreateTime),
			UpdateTime: timestamppb.New(message.UpdateTime),
		}
		if userID := message.UserID; userID != nil {
			value := formatUserName(*userID)
			messages[i].User = &value
		}
	}

	return &v1.BatchGetMessagesResponse{
		Messages: messages,
	}, nil
}

func (s *API) CreateMessage(ctx context.Context, req *v1.CreateMessageRequest) (*v1.Message, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, err := parseConversationName(req.Parent)
	if err != nil {
		return nil, err
	}

	result, err := command.Send(ctx, s.commandBus, command.CreateMessage{
		ActorID:        subject,
		ConversationID: conversationID,
		Content:        req.Message.Content,
	})
	if err != nil {
		return nil, err
	}

	m := &v1.Message{
		Name:       formatMessageName(conversationID, result.Message.ID),
		Content:    result.Message.Content,
		CreateTime: timestamppb.New(result.Message.CreateTime),
		UpdateTime: timestamppb.New(result.Message.UpdateTime),
	}
	if userID := result.Message.UserID; userID != nil {
		value := formatUserName(*userID)
		m.User = &value
	}

	return m, nil
}

func (s *API) UpdateMessage(ctx context.Context, req *v1.UpdateMessageRequest) (*v1.Message, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, messageID, err := parseMessageName(req.Message.Name)
	if err != nil {
		return nil, err
	}

	var content *string
	if req.UpdateMask == nil || req.UpdateMask.Paths == nil {
		if req.Message.Content != "" {
			value := req.Message.Content
			content = &value
		}
	} else {
		for _, path := range req.UpdateMask.Paths {
			switch path {
			case "*":
				value := req.Message.Content
				content = &value

			case "content":
				value := req.Message.Content
				content = &value

			default:
				return nil, ErrInvalidUpdateMaskPath
			}
		}
	}

	if content == nil {
		result, err := query.Send(ctx, s.queryBus, query.GetMessage{
			ActorID:        subject,
			ConversationID: conversationID,
			MessageID:      messageID,
		})
		if err != nil {
			return nil, err
		}

		m := &v1.Message{
			Name:       formatMessageName(conversationID, result.Message.ID),
			Content:    result.Message.Content,
			CreateTime: timestamppb.New(result.Message.CreateTime),
			UpdateTime: timestamppb.New(result.Message.UpdateTime),
		}
		if userID := result.Message.UserID; userID != nil {
			value := formatUserName(*userID)
			m.User = &value
		}

		return m, nil
	}

	result, err := command.Send(ctx, s.commandBus, command.EditMessage{
		ActorID:        subject,
		ConversationID: conversationID,
		MessageID:      messageID,
		Content:        *content,
	})
	if err != nil {
		return nil, err
	}

	m := &v1.Message{
		Name:       formatMessageName(conversationID, result.Message.ID),
		Content:    result.Message.Content,
		CreateTime: timestamppb.New(result.Message.CreateTime),
		UpdateTime: timestamppb.New(result.Message.UpdateTime),
	}
	if userID := result.Message.UserID; userID != nil {
		value := formatUserName(*userID)
		m.User = &value
	}

	return m, nil
}

func (s *API) DeleteMessage(ctx context.Context, req *v1.DeleteMessageRequest) (*emptypb.Empty, error) {
	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	conversationID, messageID, err := parseMessageName(req.Name)
	if err != nil {
		return nil, err
	}

	if _, err := command.Send(ctx, s.commandBus, command.DeleteMessage{
		ActorID:        subject,
		ConversationID: conversationID,
		MessageID:      messageID,
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func formatConversationUserRoleAPIV1(format string) v1.ConversationUser_Role {
	switch format {
	case "OWNER":
		return v1.ConversationUser_OWNER
	case "ADMIN":
		return v1.ConversationUser_ADMIN
	case "MEMBER":
		return v1.ConversationUser_MEMBER
	default:
		return v1.ConversationUser_ROLE_UNSPECIFIED
	}
}

func parseSubject(ctx context.Context) (string, error) {
	claims, ok := kjwt.FromContext(ctx)
	if !ok {
		return "", ErrUnauthenticated
	}

	subject, err := claims.GetSubject()
	if err != nil || subject == "" {
		return "", ErrUnauthenticated
	}

	return subject, nil
}

func parseConversationName(name string) (string, error) {
	var conversationID string
	if err := resourcename.Sscan(name, conversationResourcePattern, &conversationID); err != nil {
		return "", ErrInvalidConversationResourceName
	}

	return conversationID, nil
}

func formatConversationName(conversationID string) string {
	return resourcename.Sprint(conversationResourcePattern, conversationID)
}

func parseUserName(name string) (string, error) {
	var userID string
	if err := resourcename.Sscan(name, userResourcePattern, &userID); err != nil {
		return "", ErrInvalidUserResourceName
	}

	return userID, nil
}

func formatUserName(userID string) string {
	return resourcename.Sprint(userResourcePattern, userID)
}

func parseConversationUserName(name string) (string, string, error) {
	var conversationID, conversationUserID string
	if err := resourcename.Sscan(name, conversationUserResourcePattern, &conversationID, &conversationUserID); err != nil {
		return "", "", ErrInvalidConversationUserResourceName
	}

	return conversationID, conversationUserID, nil
}

func formatConversationUserName(conversationID string, conversationUserID string) string {
	return resourcename.Sprint(conversationUserResourcePattern, conversationID, conversationUserID)
}

func parseMessageName(name string) (string, string, error) {
	var conversationID, messageID string
	if err := resourcename.Sscan(name, messageResourcePattern, &conversationID, &messageID); err != nil {
		return "", "", ErrInvalidMessageResourceName
	}

	return conversationID, messageID, nil
}

func formatMessageName(conversationID string, messageID string) string {
	return resourcename.Sprint(messageResourcePattern, conversationID, messageID)
}
