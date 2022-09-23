package rpc

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Luoxin/sexy/us/ext"
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

	// upgrader := websocket.FastHTTPUpgrader{
	// 	HandshakeTimeout:  0,
	// 	ReadBufferSize:    0,
	// 	WriteBufferSize:   0,
	// 	WriteBufferPool:   nil,
	// 	Subprotocols:      nil,
	// 	Error:             nil,
	// 	CheckOrigin:       nil,
	// 	EnableCompression: false,
	// }

	RegisterHandler(ApiHandler{All, "/ctx", func(ctx *Ctx) error {
		log.Info("ctx")
		source := ctx.source

		// err := upgrader.Upgrade(source, func(conn *websocket.Conn) {
		// 	defer conn.Close()
		// 	for i := 0; i < 10; i++ {
		// 		time.Sleep(time.Second)
		// 		log.Info(i)
		// 		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().String()))
		// 	}
		// })
		// if err != nil {
		// 	log.Errorf("err:%v", err)
		// 	return err
		// }

		source.Response.Header.Set("Content-Type", "text/event-stream")
		source.Response.Header.Set("Cache-Control", "no-cache")
		source.Response.Header.Set("Connection", "keep-alive")
		source.Response.Header.Set("Transfer-Encoding", "chunked")
		source.Response.Header.Set("X-Content-Type-Options", "nosniff")

		// source.Response.Header.Set("Content-Type", "application/json")
		// source.Response.Header.Set("Transfer-Encoding", "chunked")
		source.Response.SetStatusCode(200)
		source.Response.ImmediateHeaderFlush = true
		source.Response.Header.SetContentLength(-1)

		// source.SetBodyStreamWriter(func(w *bufio.Writer) {
		source.Conn().Write([]byte("HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/plain\r\n" +
			"Transfer-Encoding: chunked\r\n\r\n"))

		if x, ok := source.Conn().(http.Flusher); ok {
			log.Info("flush")
			x.Flush()
		}

		// source.Conn().Write([]byte("HTTP/1.1 200 OK\\r\\r\\nAccess-Control-Allow-Origin: *\\r\\r\\nCache-Control: no-cache, no-store\r\r\nContent-Type: text/json\\r\\r\\nP3P: CP=“NOI DEV PSA PSD IVA PVD OTP OUR OTR IND OTC”\\r\\r\\nPragma: no-cache\\r\\r\\nConnection: close\\r\\r\\nTransfer-Encoding: chunked"))

		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			log.Info(i)
			// fmt.Fprintf(w, "Data:%s\n", time.Now().String())
			source.Conn().Write([]byte(time.Now().String() + "\\r\\n"))
			if x, ok := source.Conn().(http.Flusher); ok {
				log.Info("flush")
				x.Flush()
			}
			// w.Flush()
		}

		source.Conn().Write([]byte("\\r\\n"))
		if x, ok := source.Conn().(http.Flusher); ok {
			log.Info("flush")
			x.Flush()
		}
		// w.Flush()
		// })

		return nil
	}})

	go func() {
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
	}()

	err := Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
