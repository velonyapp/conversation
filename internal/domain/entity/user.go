package entity

import (
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type User struct {
	id vo.UserID
}

func NewUser(
	id vo.UserID,
) *User {
	return &User{
		id: id,
	}
}

func ReconstituteUser(
	id vo.UserID,
) *User {
	return &User{
		id: id,
	}
}

func (a *User) ID() vo.UserID {
	return a.id
}
