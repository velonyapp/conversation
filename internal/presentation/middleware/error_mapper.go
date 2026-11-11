package middleware

import (
	"context"
	"errors"

	applicationcommon "github.com/velonyapp/conversation/internal/application/common"
	domainentity "github.com/velonyapp/conversation/internal/domain/entity"
	domainrepo "github.com/velonyapp/conversation/internal/domain/repo"
	domainservice "github.com/velonyapp/conversation/internal/domain/service"
	domainvo "github.com/velonyapp/conversation/internal/domain/vo"
	presentationservice "github.com/velonyapp/conversation/internal/presentation/service"

	kerrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/middleware"
)

type ErrorMapper middleware.Middleware

func NewErrorMapper() ErrorMapper {
	return ErrorMapper(
		func(next middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req any) (any, error) {
				reply, err := next(ctx, req)
				if err == nil {
					return reply, nil
				}

				var kratosErr *kerrors.Error
				if errors.As(err, &kratosErr) && kratosErr.Code < 500 {
					return reply, err
				}

				switch {
				// Presentation Service
				case errors.Is(err, presentationservice.ErrInvalidConversationResourceName):
					return reply, kerrors.BadRequest("",
						presentationservice.ErrInvalidConversationResourceName.Error(),
					)
				case errors.Is(err, presentationservice.ErrInvalidConversationBatchSize):
					return reply, kerrors.BadRequest("",
						presentationservice.ErrInvalidConversationBatchSize.Error(),
					)
				case errors.Is(err, presentationservice.ErrInvalidConversationUserResourceName):
					return reply, kerrors.BadRequest("",
						presentationservice.ErrInvalidConversationUserResourceName.Error(),
					)
				case errors.Is(err, presentationservice.ErrInvalidConversationUserBatchSize):
					return reply, kerrors.BadRequest("",
						presentationservice.ErrInvalidConversationUserBatchSize.Error(),
					)
				case errors.Is(err, presentationservice.ErrInvalidMessageResourceName):
					return reply, kerrors.BadRequest("",
						presentationservice.ErrInvalidMessageResourceName.Error(),
					)
				case errors.Is(err, presentationservice.ErrInvalidMessageBatchSize):
					return reply, kerrors.BadRequest("",
						presentationservice.ErrInvalidMessageBatchSize.Error(),
					)
				case errors.Is(err, presentationservice.ErrInvalidUserResourceName):
					return reply, kerrors.BadRequest("",
						presentationservice.ErrInvalidUserResourceName.Error(),
					)
				case errors.Is(err, presentationservice.ErrInvalidUpdateMaskPath):
					return reply, kerrors.BadRequest("",
						presentationservice.ErrInvalidUpdateMaskPath.Error(),
					)
				case errors.Is(err, presentationservice.ErrUnauthenticated):
					return reply, kerrors.Unauthorized("",
						presentationservice.ErrUnauthenticated.Error(),
					)
				case errors.Is(err, presentationservice.ErrForbidden):
					return reply, kerrors.Forbidden("",
						presentationservice.ErrForbidden.Error(),
					)

				// Application Common
				case errors.Is(err, applicationcommon.ErrConversationNotFound):
					return reply, kerrors.NotFound("",
						applicationcommon.ErrConversationNotFound.Error(),
					)
				case errors.Is(err, applicationcommon.ErrUserNotFound):
					return reply, kerrors.NotFound("",
						applicationcommon.ErrUserNotFound.Error(),
					)
				case errors.Is(err, applicationcommon.ErrConversationUserNotFound):
					return reply, kerrors.NotFound("",
						applicationcommon.ErrConversationUserNotFound.Error(),
					)
				case errors.Is(err, applicationcommon.ErrMessageNotFound):
					return reply, kerrors.NotFound("",
						applicationcommon.ErrMessageNotFound.Error(),
					)

				// Domain Entity
				case errors.Is(err, domainentity.ErrConversationDeleted):
					return reply, kerrors.NotFound("",
						applicationcommon.ErrConversationNotFound.Error(),
					)
				case errors.Is(err, domainentity.ErrConversationUserRemoved):
					return reply, kerrors.NotFound("",
						applicationcommon.ErrConversationUserNotFound.Error(),
					)
				case errors.Is(err, domainentity.ErrMessageDeleted):
					return reply, kerrors.NotFound("",
						applicationcommon.ErrMessageNotFound.Error(),
					)

				// Domain Repository
				case errors.Is(err, domainrepo.ErrConversationVersionConflict):
					return reply, kerrors.Conflict("",
						domainrepo.ErrConversationVersionConflict.Error(),
					)
				case errors.Is(err, domainrepo.ErrConversationUserVersionConflict):
					return reply, kerrors.Conflict("",
						domainrepo.ErrConversationUserVersionConflict.Error(),
					)
				case errors.Is(err, domainrepo.ErrMessageVersionConflict):
					return reply, kerrors.Conflict("",
						domainrepo.ErrMessageVersionConflict.Error(),
					)

				// Domain Service
				case errors.Is(err, domainservice.ErrCannotActOnSelf):
					return reply, kerrors.BadRequest("",
						domainservice.ErrCannotActOnSelf.Error(),
					)
				case errors.Is(err, domainservice.ErrConversationMismatch):
					return reply, kerrors.Forbidden("",
						domainservice.ErrConversationMismatch.Error(),
					)
				case errors.Is(err, domainservice.ErrCannotElectOther):
					return reply, kerrors.Forbidden("",
						domainservice.ErrCannotElectOther.Error(),
					)
				case errors.Is(err, domainservice.ErrCannotPromoteOther):
					return reply, kerrors.Forbidden("",
						domainservice.ErrCannotPromoteOther.Error(),
					)
				case errors.Is(err, domainservice.ErrCannotDemoteOther):
					return reply, kerrors.Forbidden("",
						domainservice.ErrCannotDemoteOther.Error(),
					)
				case errors.Is(err, domainservice.ErrCannotKickOther):
					return reply, kerrors.Forbidden("",
						domainservice.ErrCannotKickOther.Error(),
					)
				case errors.Is(err, domainservice.ErrCannotUpdateConversation):
					return reply, kerrors.Forbidden("",
						domainservice.ErrCannotUpdateConversation.Error(),
					)
				case errors.Is(err, domainservice.ErrCannotDeleteConversation):
					return reply, kerrors.Forbidden("",
						domainservice.ErrCannotDeleteConversation.Error(),
					)
				case errors.Is(err, domainservice.ErrCannotLeaveConversation):
					return reply, kerrors.BadRequest("",
						domainservice.ErrCannotLeaveConversation.Error(),
					)
				case errors.Is(err, domainservice.ErrCannotEditMessage):
					return reply, kerrors.Forbidden("",
						domainservice.ErrCannotEditMessage.Error(),
					)
				case errors.Is(err, domainservice.ErrCannotDeleteMessage):
					return reply, kerrors.Forbidden("",
						domainservice.ErrCannotDeleteMessage.Error(),
					)

				// Domain Value Object
				case errors.Is(err, domainvo.ErrConversationTitleEmpty):
					return reply, kerrors.BadRequest("",
						domainvo.ErrConversationTitleEmpty.Error(),
					)
				case errors.Is(err, domainvo.ErrConversationTitleTooLong):
					return reply, kerrors.BadRequest("",
						domainvo.ErrConversationTitleTooLong.Error(),
					)
				case errors.Is(err, domainvo.ErrConversationTitleInvalidCharacter):
					return reply, kerrors.BadRequest("",
						domainvo.ErrConversationTitleInvalidCharacter.Error(),
					)
				case errors.Is(err, domainvo.ErrConversationDescriptionTooLong):
					return reply, kerrors.BadRequest("",
						domainvo.ErrConversationDescriptionTooLong.Error(),
					)
				case errors.Is(err, domainvo.ErrConversationDescriptionInvalidCharacter):
					return reply, kerrors.BadRequest("",
						domainvo.ErrConversationDescriptionInvalidCharacter.Error(),
					)
				case errors.Is(err, domainvo.ErrMessageContentTooLong):
					return reply, kerrors.BadRequest("",
						domainvo.ErrMessageContentTooLong.Error(),
					)
				case errors.Is(err, domainvo.ErrMessageContentInvalidCharacter):
					return reply, kerrors.BadRequest("",
						domainvo.ErrMessageContentInvalidCharacter.Error(),
					)

				default:
					return reply, kerrors.InternalServer("",
						"internal server error",
					).WithCause(err)
				}
			}
		},
	)
}

func (m ErrorMapper) Middleware() middleware.Middleware {
	return middleware.Middleware(m)
}
