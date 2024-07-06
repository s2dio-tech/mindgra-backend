package usecase

import (
	"context"
	"time"

	"github.com/s2dio-tech/mindgra-backend/common"
	"github.com/s2dio-tech/mindgra-backend/domain"
	"golang.org/x/exp/slog"
)

type graphUsecase struct {
	graphRepo domain.GraphRepository
}

func InitGraphUsecase(repo domain.GraphRepository) domain.GraphUsecase {
	return &graphUsecase{
		graphRepo: repo,
	}
}

func (u *graphUsecase) List(c context.Context, user domain.Profile) ([]domain.Graph, error) {
	return u.graphRepo.Select(user.Id)
}

func (u *graphUsecase) Get(c context.Context, id string) (*domain.Graph, error) {
	return u.graphRepo.SelectOne(id)
}

func (u *graphUsecase) Create(c context.Context, graph domain.Graph, user domain.Profile) (res *string, err error) {
	// insert to db
	w := &graph
	w.CreatedAt = common.ToPointer(time.Now())
	wId, err := u.graphRepo.Store(*w)
	if err != nil {
		slog.Error("Create error", err)
		return nil, common.ErrInternalServerError
	}

	return wId, nil
}

func (u *graphUsecase) Update(c context.Context, id string, graph domain.Graph, user domain.Profile) error {
	sp, err := u.graphRepo.SelectOne(id)
	if err != nil {
		return common.ErrInternalServerError
	}

	if sp == nil || (user.Role == domain.RoleMember && user.Id != sp.UserId) {
		return common.ErrNotFound
	}

	err = u.graphRepo.Update(id, graph)
	if err != nil {
		slog.Error("Update graph error", err)
	}
	return err
}

func (u *graphUsecase) Delete(c context.Context, id string, user domain.Profile) error {
	graph, err := u.graphRepo.SelectOne(id)
	if err != nil {
		return common.ErrInternalServerError
	}

	if graph == nil || (user.Role == domain.RoleMember && user.Id != graph.UserId) {
		return common.ErrNotFound
	}

	return u.graphRepo.Delete(id)

}
