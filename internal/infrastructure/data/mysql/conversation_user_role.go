package mysql

const (
	conversationUserRoleOwner  uint8 = 1
	conversationUserRoleAdmin  uint8 = 2
	conversationUserRoleMember uint8 = 3
)

func parseConversationUserRole(role uint8) string {
	switch role {
	case conversationUserRoleOwner:
		return "OWNER"
	case conversationUserRoleAdmin:
		return "ADMIN"
	case conversationUserRoleMember:
		return "MEMBER"
	default:
		return ""
	}
}

func formatConversationUserRole(role string) uint8 {
	switch role {
	case "OWNER":
		return conversationUserRoleOwner
	case "ADMIN":
		return conversationUserRoleAdmin
	case "MEMBER":
		return conversationUserRoleMember
	default:
		return 0
	}
}
