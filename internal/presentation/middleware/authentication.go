package middleware

import (
	"github.com/velonyapp/conversation/internal/conf"

	kjwt "github.com/go-kratos/kratos/contrib/middleware/jwt/v3"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/golang-jwt/jwt/v5"
)

type Authentication middleware.Middleware

func NewAuthentication(c *conf.Auth) (Authentication, error) {
	auth := kjwt.Server(
		func(_ *jwt.Token) (any, error) {
			return []byte(c.AccessToken.Secret), nil
		},
		kjwt.WithParserOptions(
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithExpirationRequired(),
		),
	)

	return Authentication(auth), nil
}

func (m Authentication) Middleware() middleware.Middleware {
	return middleware.Middleware(m)
}