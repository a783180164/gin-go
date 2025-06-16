package weather

import (
	"fmt"
	"gin-go/pkg/code"
	Core "gin-go/pkg/internal/core"
	// "gin-go/pkg/logger"
	"github.com/gin-gonic/gin"
	// "github.com/sirupsen/logrus"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type NowRequst struct {
	Location string `form:"location"  binding:"required"`
	Unit     string `form:"unit,default=c" `
	Language string `form:"language,default=zh-Hans"`
}

type result struct {
	Args    string            `json:"args"`
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
	Url     string            `json:"url"`
}

type resultErrorRespone struct {
	Status_Code string `json:"status_code"`
	Status      string `json:"status"`
}

func (h *handler) Now(c *gin.Context) {
	req := new(NowRequst)
	core := Core.NewContext(c)
	if err := core.ShouldBindQuery(req); err != nil {
		fmt.Println("err", err.Error())
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.ParamsError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	weatherReq, err := http.NewRequest("GET", "https://api.seniverse.com/v3/weather/now.json", nil)

	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.ParamsError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	params := weatherReq.URL.Query()
	params.Add("location", req.Location)
	params.Add("unit", req.Unit)
	params.Add("language", req.Language)
	params.Add("key", "S2IiZWW4U4Zq-9jrx")
	weatherReq.URL.RawQuery = params.Encode()

	weatherRes, _ := http.DefaultClient.Do(weatherReq)
	defer weatherRes.Body.Close()
	fmt.Println(weatherReq)

	body, err := ioutil.ReadAll(weatherRes.Body)

	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.JsonBodyError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	var data interface{}

	json.Unmarshal(body, &data)
	if weatherRes.StatusCode != 200 {

		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.ParamsError,
			Data:    data,
			Message: "",
		})
		return
	}
	c.JSON(http.StatusOK, &code.Failure{
		Code:    code.Success,
		Data:    data,
		Message: "",
	})
}
