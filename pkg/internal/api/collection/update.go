package collection

import (
	"gin-go/pkg/code"
	Core "gin-go/pkg/internal/core"
	"gin-go/pkg/internal/service/collection"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdateRequest struct {
	UUID   string `json:"uuid" binding:"required"`
	Prompt string `json:"prompt"`
	Desc   string `json:"desc"`
}

func (h *handler) Update(c *gin.Context) {
	core := Core.NewContext(c)
	req := new(UpdateRequest)
	if err := core.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.ParamsError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	model := new(collection.UpdateCollection)

	model.UUID = req.UUID
	model.Description = req.Desc
	model.Prompt = req.Prompt

	err := h.collectionService.Update(model)
	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.CreateCollectionError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &code.Failure{
		Code:    code.Success,
		Data:    nil,
		Message: code.Text(code.Success),
	})
}
