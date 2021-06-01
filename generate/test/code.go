package test

type Code int

//go:generate changeToolForCode "./code.go" "./language_pkg.yaml" resp
const (
	FAILED           Code = -1     // 失败
	SUCCESS          Code = 0      // 成功
	ReqArgsErr       Code = 100400 // 请求参数错误信息
	SaveRecordFailed Code = 100500 // 后台存储记录失败
)
