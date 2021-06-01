package main

type ErrCode int

//go:generate stringer -type ErrCode -output code_string.go -linecomment
const (
	ERR_CODE_OK             ErrCode = 0 // 0 - OK
	ERR_CODE_INVALID_PARAMS ErrCode = 1 // 1 - 无效参数
	ERR_CODE_TIMEOUT        ErrCode = 2 // 2 - 超时
)
