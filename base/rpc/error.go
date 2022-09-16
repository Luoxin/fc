package rpc

import (
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
)

const (
	/*
		2000 - 3000: http 状态错误
			-2000 - http status code
	*/

	Success = 0
)

type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (p *Error) Error() string {
	return fmt.Sprintf("err_code:%d,err_msg:%s", p.Code, p.Message)
}

func (p *Error) Clone() *Error {
	return &Error{
		Code:    p.Code,
		Message: p.Message,
	}
}

func (p *Error) Marshal() []byte {
	buf, err := sonic.Marshal(p)
	if err != nil {
		log.Errorf("err:%v", err)
	}
	return buf
}

var errMap = map[int32]*Error{}

func RegisterErrCode(code int32, format string, a ...any) {
	_, ok := errMap[code]
	if ok {
		log.Warnf("code dup")
	}

	errMap[code] = &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}

func CreateError(code int32) *Error {
	err, ok := errMap[code]
	if ok {
		return err.Clone()
	}

	return &Error{
		Code: code,
	}
}

func CreateErrorWithMsg(code int32, format string, a ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}
