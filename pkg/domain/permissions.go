package domain

type Permission string

const (
	PermissionAdmin       Permission = "admin"
	PermissionBroadcaster Permission = "broadcaster"
	PermissionModerator   Permission = "moderator"
	PermissionVIP         Permission = "vip"
)
