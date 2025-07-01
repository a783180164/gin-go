package ollamatest

import (
	"gin-go/pkg/code"
	Core "gin-go/pkg/internal/core"
	"gin-go/pkg/internal/service/ollamatest"
	"gin-go/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PromptRequest struct {
	Text string `json:"text" binding:"required"`
	UUID string `json:"name" binding:"uuid"`
}

func (h *handler) Prompt(c *gin.Context) {
	req := new(PromptRequest)

	core := Core.NewContext(c)
	if err := core.ShouldBindJSON(req); err != nil {
		logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Info("参数绑定错误")
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.ParamsError,
			Data:    nil,
			Message: code.Text(code.ParamsError),
		})
		return
	}
	data := new(ollamatest.Prompt)
	data.Text = req.Text

	data.UUID = req.UUID
	datas, err := h.ollamatestService.Prompt(data)

	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.ParamsError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &code.Failure{
		Code:    code.Success,
		Data:    datas,
		Message: code.Text(code.Success),
	})
}
