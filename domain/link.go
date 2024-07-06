package domain

import (
	"context"
	"time"
)

type Link struct {
	Id          string
	Word1Id     string
	Word2Id     string
	UserId      string
	Content     string
	Description *string
	Refs        *[]string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type LinkRepository interface {
	FindById(id string) (*Link, error)
	FindByWordIds(w1Id string, w2Id string) (*Link, error)
	Store(r Link) (*string, error)
	Update(id string, link Link) error
	// Delete(id string) error
	Delete(w1Id string, w2Id string) error
}

type LinkUsecase interface {
	GetDetail(id string) (*Link, error)
	GetDetailByWordIds(w1id string, w2id string) (*Link, error)
	Create(c context.Context, w1Id string, w2Id string, r Link, user Profile) (res *string, err error)
	Update(c context.Context, id string, link Link, user Profile) error
	Delete(c context.Context, w1Id string, w2Id string, user Profile) error
}
