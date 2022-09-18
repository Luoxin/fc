package rpc

import (
	"github.com/darabuchi/log"
	"github.com/valyala/fasthttp"
)

var service = NewService()

func NewService() *fasthttp.Server {
	return &fasthttp.Server{
		Handler:                            handlerRouter.Handler,
		ErrorHandler:                       nil,
		HeaderReceived:                     nil,
		ContinueHandler:                    nil,
		Name:                               "",
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
		TCPKeepalive:                       false,
		ReduceMemoryUsage:                  false,
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
		CloseOnShutdown:                    false,
		StreamRequestBody:                  false,
		ConnState:                          nil,
		Logger:                             log.Clone(),
		TLSConfig:                          nil,
	}
}

func Listen(addr string) error {
	err := service.ListenAndServe(addr)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func StartService() error {

	return nil
}
