package vo

import "errors"

var (
	ErrConversationUserRoleInvalid = errors.New("invalid conversation user role")
)

type ConversationUserRole string

const (
	ConversationUserRoleOwner  ConversationUserRole = "OWNER"
	ConversationUserRoleAdmin  ConversationUserRole = "ADMIN"
	ConversationUserRoleMember ConversationUserRole = "MEMBER"
)

func NewConversationUserRole(value string) (ConversationUserRole, error) {
	role := ConversationUserRole(value)

	switch role {
	case ConversationUserRoleOwner,
		ConversationUserRoleAdmin,
		ConversationUserRoleMember:
		return role, nil
	default:
		return "", ErrConversationUserRoleInvalid
	}
}

func (r ConversationUserRole) String() string {
	return string(r)
}

func (r ConversationUserRole) Equal(other ConversationUserRole) bool {
	return r == other
}

func (r ConversationUserRole) IsOwner() bool {
	return r == ConversationUserRoleOwner
}

func (r ConversationUserRole) IsAdmin() bool {
	return r == ConversationUserRoleAdmin
}

func (r ConversationUserRole) IsMember() bool {
	return r == ConversationUserRoleMember
}
