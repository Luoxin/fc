package rpc

import (
	"errors"
	"testing"
	"time"

	"github.com/Luoxin/fc/base/ext"
	"github.com/darabuchi/log"
)

func TestListen(t *testing.T) {
	type (
		HiReq struct {
			Hi string `json:"hi"`
		}

		HiRsp struct {
			Hi string `json:"hi"`
		}
	)

	RegisterHandler(GET, "/", func(ctx *Ctx, req *ext.ExtReq) (*ext.ExtRsp, error) {
		var rsp ext.ExtRsp

		log.Info("log")

		return &rsp, errors.New("is error")
	})
	RegisterHandler(POST, "/hi", func(ctx *Ctx, req *HiReq) (*HiRsp, error) {
		var rsp HiRsp

		log.Info(req.Hi)
		rsp.Hi = req.Hi

		return &rsp, nil
	})
	RegisterHandler(All, "/p1/p2/p3/p4", func(ctx *Ctx, req *ext.ExtReq) (*ext.ExtRsp, error) {
		var rsp ext.ExtRsp

		log.Info("log")

		rsp.Buf = []byte("p4")

		return &rsp, nil
	})
	RegisterHandler(All, "/p1/p2/p3/{p4}", func(ctx *Ctx, req *ext.ExtReq) (*ext.ExtRsp, error) {
		var rsp ext.ExtRsp

		log.Info("log")

		return &rsp, nil
	})

	go func() {
		time.Sleep(time.Second * 3)
		{
			var rsp HiRsp
			err := Call("", "/hi", &HiReq{
				Hi: "hi",
			}, &rsp)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
			log.Info(rsp.Hi)
		}

		{
			var rsp ext.ExtRsp
			err := Call("", "/", &ext.ExtReq{}, &rsp)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
			log.Info(string(rsp.Buf))
		}

		{
			var rsp ext.ExtRsp
			err := Call("", "/p1/p2/p3/p4", &ext.ExtReq{}, &rsp)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
			log.Info(string(rsp.Buf))
		}
	}()

	err := Listen(":8080")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
