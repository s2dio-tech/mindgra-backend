package domain

import (
	"context"
	"time"
)

type Graph struct {
	Id        string     `json:"id"`
	UserId    string     `json:"userId"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type GraphRepository interface {
	Select(userId string) ([]Graph, error)
	SelectOne(id string) (*Graph, error)
	Store(r Graph) (*string, error)
	Update(id string, graph Graph) error
	Delete(id string) error
}

type GraphUsecase interface {
	List(c context.Context, user Profile) ([]Graph, error)
	Get(c context.Context, id string) (*Graph, error)
	Create(c context.Context, graph Graph, user Profile) (res *string, err error)
	Update(c context.Context, id string, graph Graph, user Profile) error
	Delete(c context.Context, id string, user Profile) error
}
