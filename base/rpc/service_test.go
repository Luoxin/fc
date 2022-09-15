package rpc

import (
	"errors"
	"testing"

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

	RegisterHandler(GET, "/", func(ctx *Ctx, req *ExtReq) (*ExtRsp, error) {
		var rsp ExtRsp

		log.Info("log")

		return &rsp, errors.New("is error")
	})
	RegisterHandler(POST, "/hi", func(ctx *Ctx, req *HiReq) (*HiRsp, error) {
		var rsp HiRsp

		log.Info(req.Hi)
		rsp.Hi = req.Hi

		return &rsp, nil
	})
	RegisterHandler(All, "/p1/p2/p3/p4", func(ctx *Ctx, req *ExtReq) (*ExtRsp, error) {
		var rsp ExtRsp

		log.Info("log")

		return &rsp, nil
	})
	RegisterHandler(All, "/p1/p2/p3/{p4}", func(ctx *Ctx, req *ExtReq) (*ExtRsp, error) {
		var rsp ExtRsp

		log.Info("log")

		return &rsp, nil
	})
	Listen(":8080")
}
