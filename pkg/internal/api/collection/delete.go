package collection

import (
	"gin-go/pkg/code"
	Core "gin-go/pkg/internal/core"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type DeleteRequest struct {
	UUIDS string `json:"uuids" binding:"required"`
}

func (h *handler) Delete(c *gin.Context) {
	core := Core.NewContext(c)
	req := new(DeleteRequest)
	if err := core.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.ParamsError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	uuids := strings.Split(req.UUIDS, ",")

	err := h.collectionService.Delete(uuids)
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
