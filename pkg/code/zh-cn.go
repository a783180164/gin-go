package code

var zhCNText = map[int]string{
	Success:           "成功",
	CreateUserParam:   "参数不正确",
	CreatHasUserError: "账号已经存在",
	CreatUserNoError:  "创建账号失败",

	CreateCollectionError:     "创建知识库失败",
	CreateHaveCollection:      "知识库已经存在",
	QueryCollectionsError:     "查询列表失败",
	QueryCollectionCountError: "查询总数失败",

	HashPasswordError:  "加密失败",
	AuthorizationNo:    "没有Authorization",
	AuthorizationError: "Authorization解析失败",
	ParamsError:        "参数不正确",

	JsonBodyError: "json解析失败",
}
