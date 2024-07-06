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

type WordHandler struct {
	wordUsecase domain.WordUsecase
}

func InitWordHandlers(wus domain.WordUsecase) *WordHandler {
	return &WordHandler{
		wordUsecase: wus,
	}
}

func (h *WordHandler) CreateWord(c *gin.Context) {
	var schema WordCreateRequestSchema
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

	_word := domain.Word{
		UserId:      authCommon.ExtractUser(c).Id,
		Content:     schema.Content,
		Description: schema.Description,
		Refs:        schema.Refs,
	}

	var id *string
	var err error
	var link *domain.WordsLink = nil
	user := authCommon.ExtractUser(c)

	if schema.SourceId == nil {
		id, err = h.wordUsecase.Create(c, _word, schema.GraphId, user)
	} else {
		id, err = h.wordUsecase.CreateWordWithLink(c, _word, *schema.SourceId, schema.GraphId, user)
		if err == nil {
			link = &domain.WordsLink{
				SourceId: *schema.SourceId,
				TargetId: *id,
			}
		}
	}
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"word": domain.Word{
			Id:          *id,
			GraphId:     schema.GraphId,
			UserId:      user.Id,
			Content:     schema.Content,
			Description: schema.Description,
			Refs:        schema.Refs,
		},
		"link": link,
	})
}

func (h *WordHandler) UpdateWord(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	var schema WordUpdateRequestSchema
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

	err := h.wordUsecase.Update(c, id, domain.Word{
		Content:     schema.Content,
		Description: schema.Description,
		Refs:        schema.Refs,
	})
	if err != nil {
		httpCommon.ErrorResponse(c, common.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *WordHandler) DeleteWord(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	h.wordUsecase.Delete(c, id, authCommon.ExtractUser(c))
}

func (h *WordHandler) GetWordDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}
	w, err := h.wordUsecase.GetWordById(c, id)
	if err != nil {
		httpCommon.ErrorResponse(c, common.ErrInternalServerError)
		return
	}
	if w == nil {
		httpCommon.ErrorResponse(c, common.ErrNotFound)
		return
	}
	c.JSON(http.StatusOK, w)
}

func (h *WordHandler) GetGraphData(c *gin.Context) {
	var id = c.Param("id")
	if id == "" {
		httpCommon.ErrorResponse(c, common.ErrNotFound)
		return
	}

	data, err := h.wordUsecase.GetGraphData(c, id)
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *WordHandler) SearchWord(c *gin.Context) {
	var text = c.Query("search")
	var graphId = c.Query("graphId")
	if graphId == "" || len(text) < 3 {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}
	data, err := h.wordUsecase.SearchWord(c, text)
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *WordHandler) FindPath(c *gin.Context) {
	var fromWordId = c.Query("fromId")
	var toWordId = c.Query("toId")
	if fromWordId == "" || toWordId == "" {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}
	words, links, err := h.wordUsecase.FindPath(c, fromWordId, toWordId)
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"words": words,
		"links": links,
	})
}

func (h *WordHandler) Link2Words(c *gin.Context) {
	var schema Link2WordsRequestSchema
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

	err := h.wordUsecase.Link2Words(c, schema.SourceId, schema.TargetId)
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
