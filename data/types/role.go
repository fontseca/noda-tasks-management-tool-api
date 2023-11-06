package types

type Role uint8

const (
	RoleAdmin Role = 1 + iota
	RoleUser
)
