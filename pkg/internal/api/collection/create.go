package collection

import (
	"gin-go/pkg/code"
	Core "gin-go/pkg/internal/core"
	"gin-go/pkg/internal/service/collection"
	"gin-go/pkg/internal/service/ollamatest"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CollectionRequest struct {
	Name   string `json:"name" binding:"required"`
	Prompt string `json:"prompt"`
	Desc   string `json:"desc"`
}

func (h *handler) Create(c *gin.Context) {
	core := Core.NewContext(c)
	req := new(CollectionRequest)
	if err := core.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.ParamsError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	model := new(collection.CreateCollection)

	model.Name = req.Name
	model.Description = req.Desc
	model.Prompt = req.Prompt

	uuid, err := h.collectionService.Create(model)
	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.CreateCollectionError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	qdModel := new(ollamatest.CreateCollection)

	qdModel.Name = req.Name
	qdModel.UUID = uuid
	qdModel.Size = 2048
	qdErr := h.ollamatestService.Create(qdModel)
	if qdErr != nil {

		h.collectionService.Delete([]string{uuid})
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.CreateCollectionError,
			Data:    nil,
			Message: qdErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &code.Failure{
		Code:    code.Success,
		Data:    uuid,
		Message: code.Text(code.Success),
	})
}
