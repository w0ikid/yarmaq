package models

const (
	RoleAdmin    = "admin"
	RoleSupport  = "support"
	// without role is a customer with no special permissions, can only access their own account
)

func IsValidRole(role string) bool {
	switch role {
	case RoleAdmin, RoleSupport:
		return true
	default:
		return false
	}
}
