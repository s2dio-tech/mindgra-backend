package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/s2dio-tech/mindgra-backend/common"
	authCommon "github.com/s2dio-tech/mindgra-backend/common/auth"
	httpCommon "github.com/s2dio-tech/mindgra-backend/common/http"
	"github.com/s2dio-tech/mindgra-backend/domain"
)

type LinkHandler struct {
	linkUsecase domain.LinkUsecase
}

func InitLinkHandlers(us domain.LinkUsecase) *LinkHandler {
	return &LinkHandler{
		linkUsecase: us,
	}
}

func (h *LinkHandler) CreateLink(c *gin.Context) {
	var schema LinkCreateRequestSchema
	// bind request context to data struct
	if err := c.Bind(&schema); err != nil {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	// validator data struct
	v := validator.New()
	if err := v.Struct(schema); err != nil {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	id, err := h.linkUsecase.Create(
		c,
		schema.Word1Id,
		schema.Word2Id,
		domain.Link{
			Content:     schema.Content,
			Description: schema.Description,
			Refs:        schema.Refs,
		},
		authCommon.ExtractUser(c),
	)
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          id,
		"word1Id":     schema.Word1Id,
		"word2Id":     schema.Word2Id,
		"content":     schema.Content,
		"description": schema.Description,
		"refs":        schema.Refs,
	})
}

func (h *LinkHandler) UpdateLink(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	var schema LinkUpdateRequestSchema
	// bind request context to data struct
	if err := c.Bind(&schema); err != nil {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	// validator data struct
	v := validator.New()
	if err := v.Struct(schema); err != nil {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	h.linkUsecase.Update(
		c,
		id,
		domain.Link{
			Content:     schema.Content,
			Description: schema.Description,
			Refs:        schema.Refs,
		},
		authCommon.ExtractUser(c),
	)
}

func (h *LinkHandler) GetDetail(c *gin.Context) {
	path1 := c.Param("path1")
	path2 := c.Param("path2")

	var r *domain.Link
	var err error

	if path2 == "" {
		var id = path1

		r, err = h.linkUsecase.GetDetail(id)
		if err != nil {
			httpCommon.ErrorResponse(c, err)
			return
		}

	} else {

		r, err = h.linkUsecase.GetDetailByWordIds(path1, path2)
		if err != nil {
			httpCommon.ErrorResponse(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          r.Id,
		"userId":      r.UserId,
		"content":     r.Content,
		"description": r.Description,
		"refs":        r.Refs,
		"createdAt":   r.CreatedAt,
	})
}

func (h *LinkHandler) DeleteLink(c *gin.Context) {
	var schema LinkRemoveRequestSchema
	// bind request context to data struct
	if err := c.Bind(&schema); err != nil {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	// validator data struct
	v := validator.New()
	if err := v.Struct(schema); err != nil {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	err := h.linkUsecase.Delete(c, schema.Word1Id, schema.Word2Id, authCommon.ExtractUser(c))
	if err != nil {
		httpCommon.ErrorResponse(c, err)
	}

	c.JSON(http.StatusOK, nil)
}
