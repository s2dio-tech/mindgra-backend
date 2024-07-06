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

type GraphHandler struct {
	graphUsecase domain.GraphUsecase
}

func InitGraphHandlers(sus domain.GraphUsecase) *GraphHandler {
	return &GraphHandler{
		graphUsecase: sus,
	}
}

func (h *GraphHandler) List(c *gin.Context) {
	res, err := h.graphUsecase.List(c, authCommon.ExtractUser(c))
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *GraphHandler) Detail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	res, err := h.graphUsecase.Get(c, id)
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *GraphHandler) CreateGraph(c *gin.Context) {
	var schema GraphCreateRequestSchema
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

	_graph := domain.Graph{
		UserId: authCommon.ExtractUser(c).Id,
		Name:   schema.Name,
	}

	var id *string
	var err error

	id, err = h.graphUsecase.Create(c, _graph, authCommon.ExtractUser(c))
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, domain.Graph{
		Id:   *id,
		Name: schema.Name,
	})
}

func (h *GraphHandler) UpdateGraph(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	// validator data struct
	var schema GraphUpdateRequestSchema
	if err := c.Bind(&schema); err != nil {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}
	if err := validator.New().Struct(schema); err != nil {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	err := h.graphUsecase.Update(c, id, domain.Graph{
		Name: schema.Name,
		Type: schema.Type,
	},
		authCommon.ExtractUser(c),
	)

	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *GraphHandler) DeleteGraph(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	err := h.graphUsecase.Delete(c, id, authCommon.ExtractUser(c))
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
