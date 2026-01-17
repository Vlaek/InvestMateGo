package models

type UserRole string

const (
	Default    UserRole = "DEFAULT"
	Subscriber UserRole = "SUBSCRIBER"
	Admin      UserRole = "ADMIN"
)

// Проверка роли на валидность
func (r UserRole) IsValid() bool {
	switch r {
	case Default, Subscriber, Admin:
		return true
	default:
		return false
	}
}

// Cписок всех допустимых ролей
func ValidRoles() []UserRole {
	return []UserRole{Default, Subscriber, Admin}
}

// Cписок всех допустимых ролей в виде строк
func ValidRoleStrings() []string {
	return []string{
		string(Default),
		string(Subscriber),
		string(Admin),
	}
}
