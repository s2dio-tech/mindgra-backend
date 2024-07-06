package domain

import (
	"context"
	"time"
)

type RoleType string

const (
	RoleAdmin     RoleType = "admin"
	RoleModerator RoleType = "moderator"
	RoleMember    RoleType = "member"
)

var RoleMap = map[string]RoleType{
	"admin":     RoleAdmin,
	"moderator": RoleModerator,
	"member":    RoleMember,
}

type User struct {
	Id        string
	Name      string
	Email     string
	Password  string
	Role      RoleType
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type Profile struct {
	Id   string
	Role RoleType
}

type UserRepository interface {
	Create(*User) (*string, error)
	Update(string, map[string]interface{}) error
	FindById(string) (*User, error)
	FindByEmail(string) (*User, error)
}

type UserUsecase interface {
	Registration(context.Context, *User) (*string, error)
	FindByEmail(context.Context, string) (*User, error)
	FindById(context.Context, string) (*User, error)
	UpdatePassword(c context.Context, userId string, password string) error
}
