package main

import (
	"time"

	"github.com/Luoxin/sexy/base/ext"
	"github.com/Luoxin/sexy/base/rpc"
	"github.com/Luoxin/sexy/honoka"
	"github.com/darabuchi/log"
)

func main() {
	err := honoka.Load()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	rpc.RegisterHandler(rpc.All, "", func(ctx *rpc.Ctx, req *ext.ExtReq) (*ext.ExtRsp, error) {
		var rsp ext.ExtRsp

		log.Info("log")

		rsp.Buf = []byte(time.Now().String())

		return &rsp, nil
	})

	err = rpc.StartService()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
