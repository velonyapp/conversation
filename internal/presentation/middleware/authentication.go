package middleware

import (
	"context"

	"github.com/velonyapp/conversation/internal/conf"

	kjwt "github.com/go-kratos/kratos/contrib/middleware/jwt/v3"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/selector"
	"github.com/golang-jwt/jwt/v5"
)

type Authentication middleware.Middleware

func NewAuthentication(c *conf.Auth) (Authentication, error) {
	public := map[string]struct{}{
		"/grpc.health.v1.Health/Check": {},
		"/grpc.health.v1.Health/Watch": {},
	}

	auth := kjwt.Server(
		func(_ *jwt.Token) (any, error) {
			return []byte(c.AccessToken.Secret), nil
		},
		kjwt.WithParserOptions(
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithExpirationRequired(),
		),
	)

	selected := middleware.Chain(
		selector.Server(auth).
			Match(func(_ context.Context, operation string) bool {
				_, ok := public[operation]
				return !ok
			}).
			Build(),
	)

	return Authentication(selected), nil
}

func (m Authentication) Middleware() middleware.Middleware {
	return middleware.Middleware(m)
}
