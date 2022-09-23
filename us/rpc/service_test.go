package rpc

import (
	"errors"
	"fmt"
	"testing"
	"time"
	
	"github.com/Luoxin/sexy/us/ext"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
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
	
	const port = 10001
	
	RegisterHandler(ApiHandler{GET, "/", func(ctx *Ctx, req *ext.ExtReq) (*ext.ExtRsp, error) {
		var rsp ext.ExtRsp
		
		log.Info("log")
		
		return &rsp, errors.New("is error")
	}})
	RegisterHandler(ApiHandler{POST, "/hi", func(ctx *Ctx, req *HiReq) (*HiRsp, error) {
		var rsp HiRsp
		
		log.Info(req.Hi)
		rsp.Hi = req.Hi
		
		return &rsp, nil
	}})
	RegisterHandler(ApiHandler{All, "/p1/p2/p3/p4", func(ctx *Ctx, req *ext.ExtReq) (*ext.ExtRsp, error) {
		var rsp ext.ExtRsp
		
		log.Info("log")
		
		rsp.Buf = []byte("p4")
		
		return &rsp, nil
	}})
	RegisterHandler(ApiHandler{All, "/p1/p2/p3/{p5}", func(ctx *Ctx, req *ext.ExtReq) (*ext.ExtRsp, error) {
		var rsp ext.ExtRsp
		
		log.Info("log")
		
		return &rsp, nil
	}})
	
	go func() {
		time.Sleep(time.Second * 3)
		{
			var rsp HiRsp
			err := CallWithAddr("", fmt.Sprintf("http://127.0.0.1:%d", port), "/hi", &HiReq{
				Hi: "hi",
			}, &rsp)
			if err != nil {
				log.Errorf("err:%v", err)
			} else {
				log.Info(rsp.Hi)
			}
		}
		
		{
			var rsp ext.ExtRsp
			err := CallWithAddr("", fmt.Sprintf("http://127.0.0.1:%d", port), "/", &ext.ExtReq{}, &rsp)
			if err != nil {
				log.Errorf("err:%v", err)
			} else {
				log.Info(string(rsp.Buf))
			}
		}
		
		{
			var rsp ext.ExtRsp
			err := CallWithAddr("", fmt.Sprintf("http://127.0.0.1:%d", port), "/p1/p2/p3/p4", &ext.ExtReq{}, &rsp)
			if err != nil {
				log.Errorf("err:%v", err)
			} else {
				log.Info(string(rsp.Buf))
			}
		}
		utils.Exit()
	}()
	
	err := Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
