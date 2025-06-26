package code

import (
	_ "embed"
)

//go:embed code.go
var ByteCodeFile []byte

// Failure 错误时返回结构
type Failure struct {
	Code    int         `json:"code"` // 业务码
	Data    interface{} `json:"data"`
	Message string      `json:"message"` // 描述信息

}

const (
	Success           = 0
	HashPasswordError = 90100

	AuthorizationNo    = 90200
	AuthorizationError = 90300

	ParamsError   = 90400
	JsonBodyError = 90500

	CreatHasUserError = 10100
	CreatUserNoError  = 10200
	CreateUserParam   = 10300

	CreateCollectionError     = 20100
	CreateHaveCollection      = 20200
	QueryCollectionsError     = 20300
	QueryCollectionCountError = 20400
)

func Text(code int) string {

	return zhCNText[code]
}
