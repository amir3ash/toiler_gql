package model

import (
	"toiler-graphql/database"
)

type UserUser struct {
	database.UserUser
}

var prefixPath string

func SetAvatarsPrefixPath(prefix string) {
	prefixPath = prefix
}

func NormalizeUsersAvatar(u []database.UserUser) []UserUser {
	users := make([]UserUser, len(u))
	for i, v := range u {
		users[i] = UserUser{v}

		if v.Avatar.Valid && v.Avatar.String != "" {
			users[i].Avatar.String = prefixPath + v.Avatar.String
		}
	}
	return users
}
