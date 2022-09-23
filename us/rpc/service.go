package rpc

import (
	"errors"
	"fmt"
	"net"
	
	"github.com/Luoxin/sexy/honoka"
	"github.com/Luoxin/sexy/us/ext"
	"github.com/Luoxin/sexy/us/nozomi"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/valyala/fasthttp"
)

var service = NewService()

func NewService() *fasthttp.Server {
	return &fasthttp.Server{
		Handler: handlerRouter.Handler,
		ErrorHandler: func(ctx *fasthttp.RequestCtx, err error) {
			log.Errorf("err:%v", err)
		},
		// HeaderReceived: func(header *fasthttp.RequestHeader) fasthttp.RequestConfig {
		// 	return fasthttp.RequestConfig{
		// 		ReadTimeout:        0,
		// 		WriteTimeout:       0,
		// 		MaxRequestBodySize: 0,
		// 	}
		// },
		ContinueHandler:                    nil,
		Name:                               honoka.ServiceName,
		Concurrency:                        0,
		ReadBufferSize:                     0,
		WriteBufferSize:                    0,
		ReadTimeout:                        0,
		WriteTimeout:                       0,
		IdleTimeout:                        0,
		MaxConnsPerIP:                      0,
		MaxRequestsPerConn:                 0,
		MaxIdleWorkerDuration:              0,
		TCPKeepalivePeriod:                 0,
		MaxRequestBodySize:                 0,
		DisableKeepalive:                   false,
		TCPKeepalive:                       true,
		ReduceMemoryUsage:                  true,
		GetOnly:                            false,
		DisablePreParseMultipartForm:       false,
		LogAllErrors:                       false,
		SecureErrorLogMessage:              false,
		DisableHeaderNamesNormalizing:      false,
		SleepWhenConcurrencyLimitsExceeded: 0,
		NoDefaultServerHeader:              false,
		NoDefaultDate:                      false,
		NoDefaultContentType:               false,
		KeepHijackedConns:                  false,
		CloseOnShutdown:                    true,
		StreamRequestBody:                  true,
		ConnState:                          nil,
		Logger:                             log.Clone(),
		TLSConfig:                          nil,
	}
}

func Listen(addr string) error {
	log.Infof("listen %s", addr)
	err := service.ListenAndServe(addr)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	
	return nil
}

func StartService() (err error) {
	honoka.BindIp, err = GetLocalIp()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	
	// 添加系统级的接口
	RegisterHandler(ApiHandler{Method: HEAD, Path: "_sys_/Echo", Handler: func(ctx *Ctx, req *ext.ExtReq) (*ext.ExtRsp, error) {
		var rsp ext.ExtRsp
		return &rsp, nil
	}})
	
	// 监听端口，注册服务
	
	addr := fmt.Sprintf("%s:%d", honoka.BindIp, honoka.ConfigGet(honoka.ListenPort))
	node := &nozomi.Node{
		Address: fmt.Sprintf("http://%s", addr),
		State:   0,
	}
	
	err = nozomi.RegisterServer(node)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer nozomi.UnregisterServer()
	
	go func() {
		err = Listen(addr)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
		defer func() {
			err = service.Shutdown()
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
		}()
	}()
	
	<-utils.GetExitSign()
	
	return nil
}

// TODO: 性能优化
func GetLocalIp() (string, error) {
	addressList, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addressList {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if utils.IsLocalIp(ipNet.IP.String()) {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("not found usable ip")
}
