package domain

type Permission string

const (
	PermAll            Permission = "*"
	PermCreateResource Permission = "resource:create"
)

var ValidPermissions = map[Permission]struct{}{
	PermAll:            {},
	PermCreateResource: {},
}
