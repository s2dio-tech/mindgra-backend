package usecase

import (
	"context"
	"time"

	"github.com/s2dio-tech/mindgra-backend/common"
	"github.com/s2dio-tech/mindgra-backend/domain"
	"golang.org/x/exp/slog"
)

type wordUsecase struct {
	wordRepo  domain.WordRepository
	graphRepo domain.GraphRepository
}

func InitWordUsecase(repo domain.WordRepository, spRepo domain.GraphRepository) domain.WordUsecase {
	return &wordUsecase{
		wordRepo:  repo,
		graphRepo: spRepo,
	}
}

func (u *wordUsecase) GetGraphData(c context.Context, graphId string) (data *domain.WordsGraphData, err error) {
	ws, ls, err := u.wordRepo.FindByGraphId(graphId)

	if err != nil {
		slog.Error("GetGraphData error", err)
		return nil, common.ErrInternalServerError
	}

	return &domain.WordsGraphData{
		Words: ws,
		Links: ls,
	}, nil
}

func (u *wordUsecase) SearchWord(c context.Context, text string) ([]domain.Word, error) {
	res, err := u.wordRepo.FindByContentOrDescription(text, 10)
	if err != nil {
		slog.Error("FindByContentOrDescription error", err)
		return nil, common.ErrInternalServerError
	}
	return res, nil
}

func (u *wordUsecase) GetWordById(c context.Context, id string) (*domain.Word, error) {
	return u.wordRepo.FindById(id)
}

func (u *wordUsecase) Create(c context.Context, word domain.Word, graphId string, user domain.Profile) (res *string, err error) {
	sp, err := u.graphRepo.SelectOne(graphId)
	if err != nil {
		return nil, common.ErrInternalServerError
	}

	if sp == nil || (user.Role == domain.RoleMember && user.Id != sp.UserId) {
		return nil, common.ErrNotFound
	}

	// insert to db
	w := &word
	w.CreatedAt = common.ToPointer(time.Now())
	wId, err := u.wordRepo.Store(*w, graphId, nil)
	if err != nil {
		slog.Error("Create error", err)
		return nil, common.ErrInternalServerError
	}

	return wId, nil
}

func (u *wordUsecase) CreateWordWithLink(c context.Context, word domain.Word, linkWordId string, graphId string, user domain.Profile) (res *string, err error) {
	sp, err := u.graphRepo.SelectOne(graphId)
	if err != nil {
		return nil, common.ErrInternalServerError
	}

	if sp == nil || (user.Role == domain.RoleMember && user.Id != sp.UserId) {
		return nil, common.ErrNotFound
	}

	// validate that link word is existed or not
	joinWord, err := u.wordRepo.FindById(linkWordId)
	if err != nil {
		return nil, common.ErrInternalServerError
	}
	if joinWord == nil {
		return nil, common.ErrNotFound
	}

	wId, err := u.wordRepo.Store(word, graphId, &linkWordId)
	if err != nil {
		return nil, common.ErrInternalServerError
	}

	return wId, nil
}

func (u *wordUsecase) Update(c context.Context, id string, word domain.Word) error {
	w, err := u.wordRepo.FindById(id)
	if err != nil {
		return common.ErrInternalServerError
	}
	if w == nil {
		return common.ErrNotFound
	}

	return u.wordRepo.Update(domain.Word{
		Id:          id,
		Content:     word.Content,
		Description: word.Description,
		Refs:        word.Refs,
	})
}

func (u *wordUsecase) Delete(c context.Context, id string, user domain.Profile) error {
	word, err := u.wordRepo.FindById(id)
	if err != nil {
		return common.ErrInternalServerError
	}

	if user.Role == domain.RoleMember && user.Id != word.UserId {
		return common.ErrUnauthorization
	}

	return u.wordRepo.Delete(id)

}

func (u *wordUsecase) FindPath(c context.Context, fromId string, toId string) ([]domain.Word, []domain.WordsLink, error) {
	return u.wordRepo.FindPath(fromId, toId)
}

func (u *wordUsecase) Link2Words(c context.Context, sourceId string, targetId string) error {
	if sourceId == targetId {
		return common.ErrBadParamInput
	}
	return u.wordRepo.StoreRelation(sourceId, targetId)
}
