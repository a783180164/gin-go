package ollamatest

import (
	"gin-go/pkg/code"
	Core "gin-go/pkg/internal/core"
	qdrantModel "gin-go/pkg/internal/repository/ollamatest"
	"gin-go/pkg/internal/service/ollamatest"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListRequest struct {
	Name     string `form:"name"`
	Prompt   string `form:"prompt"`
	Desc     string `form:"desc"`
	UUID     string `form:"uuid" binding:"required"`
	Page     int    `form:"page"`
	PageSize int    `form:"pagesize"`
}

type ListData struct {
	Id   int32  `json:"id"`   // ID
	UUID string `json:"uuid"` // 现在把它加进来了
	qdrantModel.CollectionPoint
}

type ListResponse struct {
	List       []*ListData `json:"list"`
	Pagination struct {
		Total        int `json:"total"`
		CurrentPage  int `json:"current_page"`
		PerPageCount int `json:"per_page_count"`
	} `json:"pagination"`
}

func (h *handler) List(c *gin.Context) {

	req := new(ListRequest)
	core := Core.NewContext(c)
	if err := core.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.ParamsError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	page := req.Page
	if page == 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	model := new(ollamatest.ListCollections)
	model.Page = page
	model.PageSize = pageSize
	model.UUID = req.UUID
	list, err := h.ollamatestService.List(model)

	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.Success,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	count, err := h.ollamatestService.Count(model)

	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.Success,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	res := new(ListResponse)
	res.List = make([]*ListData, len(list))
	for k, v := range list {

		payload := v.GetPayload()
		cp := qdrantModel.CollectionPoint{
			Content:    payload["content"].GetStringValue(),
			ChunkIndex: int64(payload["chunk_index"].GetDoubleValue()),
			ChunkSize:  int64(payload["chunk_size"].GetDoubleValue()),
			Filename:   payload["filename"].GetStringValue(),
			CreatedAt:  payload["created_at"].GetStringValue(),
		}
		data := &ListData{
			UUID:            v.GetId().GetUuid(),
			CollectionPoint: cp,
		}

		res.List[k] = data
	}

	res.Pagination.Total = int(count)
	res.Pagination.CurrentPage = page
	res.Pagination.PerPageCount = pageSize
	c.JSON(http.StatusOK, &code.Failure{
		Code:    code.Success,
		Data:    res,
		Message: code.Text(code.Success),
	})
}
