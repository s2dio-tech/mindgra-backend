package domain

import (
	"context"
	"time"
)

type Word struct {
	Id          string     `json:"id"`
	GraphId     string     `json:"graphId"`
	UserId      string     `json:"userId"`
	Content     string     `json:"content"`
	Description *string    `json:"description"`
	Refs        *[]string  `json:"refs"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

type WordsLink struct {
	SourceId string `json:"sourceId"`
	TargetId string `json:"targetId"`
}

type WordsGraphData struct {
	Words []Word      `json:"words"`
	Links []WordsLink `json:"links"`
}

type WordRepository interface {
	FindByIds(id []string) ([]Word, error)
	FindById(id string) (*Word, error)
	FindByRandomId() (*Word, error)
	FindByGraphId(graphId string) ([]Word, []WordsLink, error)
	FindNeighborIds(id string, depth int) ([]WordsLink, error)
	FindByContentOrDescription(search string, limit int) ([]Word, error)
	FindPath(fromId string, toId string) ([]Word, []WordsLink, error)
	Store(w Word, graphId string, linkWordId *string) (*string, error)
	Update(w Word) error
	Delete(id string) error
	StoreRelation(sourceId string, targetId string) error
}

type WordUsecase interface {
	GetGraphData(c context.Context, graphId string) (data *WordsGraphData, err error)
	SearchWord(c context.Context, search string) ([]Word, error)
	FindPath(c context.Context, fromId string, toId string) ([]Word, []WordsLink, error)
	GetWordById(c context.Context, id string) (*Word, error)
	Create(c context.Context, w Word, graphId string, user Profile) (res *string, err error)
	CreateWordWithLink(c context.Context, word Word, linkWordId string, graphId string, user Profile) (res *string, err error)
	Update(c context.Context, wordId string, data Word) (err error)
	Delete(c context.Context, id string, user Profile) error
	Link2Words(c context.Context, sourceId string, targetId string) error
}
