package collection

import (
	"gin-go/pkg/code"
	Core "gin-go/pkg/internal/core"
	"gin-go/pkg/internal/service/collection"
	"gin-go/pkg/timeutil"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

type ListRequest struct {
	Name   string `form:"name"`
	Prompt string `form:"prompt"`
	Desc   string `form:"desc"`

	Page     int `form:"page"`
	PageSize int `form:"pagesize"`
}

type listData struct {
	Id        int32  `json:"id"`         // ID
	UUID      string `json:"uuid"`       // hashid
	Desc      string `json:"Desc"`       // 用户名
	Prompt    string `json:"prompt"`     // 昵称
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}
type ListResponse struct {
	List       []*listData `json:"list"`
	Pagination struct {
		Total        int `json:"total"`
		CurrentPage  int `json:"current_page"`
		PerPageCount int `json:"per_page_count"`
	} `json:"pagination"`
}

func (h *handler) List(c *gin.Context) {
	core := Core.NewContext(c)
	req := new(ListRequest)
	res := new(ListResponse)
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

	model := new(collection.ListCollection)

	model.Name = req.Name
	model.Description = req.Desc
	model.Prompt = req.Prompt

	model.Page = page
	model.PageSize = pageSize

	list, err := h.collectionService.List(model)
	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.QueryCollectionsError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	total, err := h.collectionService.Count(model)

	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.QueryCollectionCountError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	res.Pagination.Total = cast.ToInt(total)
	res.Pagination.PerPageCount = pageSize
	res.Pagination.CurrentPage = page
	res.List = make([]*listData, len(list))
	for k, v := range list {
		data := &listData{
			Id:        v.ID,
			UUID:      v.UUID,
			Desc:      v.Description,
			Prompt:    v.Prompt,
			CreatedAt: v.CreatedAt.Format(timeutil.CSTLayout),
			UpdatedAt: v.UpdatedAt.Format(timeutil.CSTLayout),
		}
		res.List[k] = data
	}

	c.JSON(http.StatusOK, &code.Failure{
		Code:    code.Success,
		Data:    res,
		Message: code.Text(code.Success),
	})
}
