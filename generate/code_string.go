// Code generated by "stringer -type ErrCode -output code_string.go -linecomment"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ERR_CODE_OK-0]
	_ = x[ERR_CODE_INVALID_PARAMS-1]
	_ = x[ERR_CODE_TIMEOUT-2]
}

const _ErrCode_name = "0 - OK1 - 无效参数2 - 超时"

var _ErrCode_index = [...]uint8{0, 6, 22, 32}

func (i ErrCode) String() string {
	if i < 0 || i >= ErrCode(len(_ErrCode_index)-1) {
		return "ErrCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrCode_name[_ErrCode_index[i]:_ErrCode_index[i+1]]
}
