package usecase

import (
	"context"

	"github.com/s2dio-tech/mindgra-backend/common"
	"github.com/s2dio-tech/mindgra-backend/domain"
)

type linkUsecase struct {
	linkRepo domain.LinkRepository
	wordRepo domain.WordRepository
}

func InitLinkUsecase(repo domain.LinkRepository, wordRepo domain.WordRepository) domain.LinkUsecase {
	return &linkUsecase{
		linkRepo: repo,
		wordRepo: wordRepo,
	}
}

func (u *linkUsecase) Create(c context.Context, w1Id string, w2Id string, link domain.Link, user domain.Profile) (*string, error) {
	// validate that link word is existed or not
	w1, err1 := u.wordRepo.FindById(w1Id)
	w2, err2 := u.wordRepo.FindById(w2Id)
	if err1 != nil || err2 != nil {
		return nil, common.ErrInternalServerError
	}
	if w1 == nil || w2 == nil {
		return nil, common.ErrBadParamInput
	}

	// if link of two words is existed
	// just update
	r, err := u.linkRepo.FindByWordIds(w1Id, w2Id)
	if err != nil {
		return nil, common.ErrInternalServerError
	}
	if r != nil {
		u.linkRepo.Update(r.Id, domain.Link{
			Content:     link.Content,
			Description: link.Description,
			Refs:        link.Refs,
		})
		return &r.Id, nil
	}

	// or not, create new
	id, err := u.linkRepo.Store(domain.Link{
		UserId:      user.Id,
		Word1Id:     w1Id,
		Word2Id:     w2Id,
		Content:     link.Content,
		Description: link.Description,
		Refs:        link.Refs,
	})
	if err != nil {
		return nil, common.ErrInternalServerError
	}
	return id, nil
}

func (u *linkUsecase) Update(c context.Context, id string, link domain.Link, user domain.Profile) error {
	r, err := u.linkRepo.FindById(id)
	if err != nil {
		return common.ErrInternalServerError
	}
	if r == nil {
		return common.ErrNotFound
	}
	if r.UserId != user.Id && user.Role == domain.RoleMember {
		return common.ErrUnauthorization
	}

	return u.linkRepo.Update(id, link)
}

func (u *linkUsecase) GetDetail(id string) (*domain.Link, error) {
	r, err := u.linkRepo.FindById(id)
	if err != nil {
		return nil, common.ErrInternalServerError
	}
	if r == nil {
		return nil, common.ErrNotFound
	}
	return r, nil
}

func (u *linkUsecase) GetDetailByWordIds(w1id string, w2id string) (*domain.Link, error) {
	w1, err1 := u.wordRepo.FindById(w1id)
	w2, err2 := u.wordRepo.FindById(w2id)
	if err1 != nil || err2 != nil {
		return nil, common.ErrInternalServerError
	}
	if w1 == nil || w2 == nil {
		return nil, common.ErrNotFound
	}

	r, err := u.linkRepo.FindByWordIds(w1id, w2id)
	if err != nil {
		return nil, common.ErrInternalServerError
	}
	if r == nil {
		return nil, common.ErrNotFound
	}
	return r, nil
}

func (u *linkUsecase) Delete(c context.Context, w1Id string, w2Id string, user domain.Profile) error {
	ws, err := u.wordRepo.FindByIds([]string{w1Id, w2Id})
	if err != nil {
		return common.ErrInternalServerError
	}
	if len(ws) != 2 || ws[0].UserId != user.Id || ws[1].UserId != user.Id {
		return common.ErrNotFound
	}
	// if user.Role == domain.RoleMember && r.UserId != user.Id {
	// 	return common.ErrUnauthorization
	// }
	// return u.linkRepo.Delete(r.Id)
	return u.linkRepo.Delete(w1Id, w2Id)
}
