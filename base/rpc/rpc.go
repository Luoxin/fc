package rpc

import (
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
)

type Ctx struct {
}

type (
	ExtReq struct {
	}

	ExtRsp struct {
		Buf []byte `json:"buf"`
	}

	BaseRsp struct {
		Error
		Data any `json:"data,omitempty"`
	}

	Error struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
	}
)

func (p *Error) Error() string {
	return fmt.Sprintf("err_code:%d,err_msg:%s", p.Code, p.Message)
}

func (p *Error) Marshal() []byte {
	buf, err := sonic.Marshal(p)
	if err != nil {
		log.Errorf("err:%v", err)
	}
	return buf
}
