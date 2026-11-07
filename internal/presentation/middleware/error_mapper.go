package middleware

import (
	"context"
	"errors"

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
				case errors.Is(err, presentationservice.ErrUnauthenticated):
					return reply, kerrors.Unauthorized("",
						presentationservice.ErrUnauthenticated.Error(),
					)
				case errors.Is(err, presentationservice.ErrForbidden):
					return reply, kerrors.Forbidden("",
						presentationservice.ErrForbidden.Error(),
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
