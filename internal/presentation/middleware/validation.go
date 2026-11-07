package middleware

import (
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/validate"
	"go.einride.tech/aip/fieldbehavior"
	"google.golang.org/protobuf/proto"
)

type Validation middleware.Middleware

func NewValidation() Validation {
	return Validation(
		validate.Validator(func(req any) error {
			message, ok := req.(proto.Message)
			if !ok {
				return nil
			}

			return fieldbehavior.ValidateRequiredFields(message)
		}),
	)
}

func (m Validation) Middleware() middleware.Middleware {
	return middleware.Middleware(m)
}
