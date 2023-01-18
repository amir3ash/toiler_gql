package model

import (
	"database/sql"
	"toiler-graphql/database"
)

type UserUser struct {
	ID        int32
	Username  string
	FirstName string
	LastName  string
	Avatar    sql.NullString
}

var prefixPath string

func SetAvatarsPrefixPath(prefix string) {
	prefixPath = prefix
}

func NormalizeUsersAvatar(u []database.UserUser) []UserUser {
	users := make([]UserUser, len(u))
	for i, v := range u {
		users[i] = UserUser{
			ID:        v.ID,
			Username:  v.Username,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Avatar:    v.Avatar,
		}

		if v.Avatar.Valid && v.Avatar.String != "" {
			users[i].Avatar.String = prefixPath + v.Avatar.String
		}
	}
	return users
}
