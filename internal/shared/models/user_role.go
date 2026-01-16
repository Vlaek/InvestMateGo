package models

type UserRole string

const (
	Default    UserRole = "DEFAULT"
	Subscriber UserRole = "SUBSCRIBER"
	Admin      UserRole = "ADMIN"
)
