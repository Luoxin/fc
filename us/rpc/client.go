package rpc

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	
	"github.com/Luoxin/sexy/us/ext"
	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/valyala/fasthttp"
)

var client = NewClient()

func NewClient() *fasthttp.Client {
	return &fasthttp.Client{
		Name:                     "",
		NoDefaultUserAgentHeader: true,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      100,
			DNSCacheDuration: time.Hour,
		}).Dial,
		DialDualStack:                 true,
		TLSConfig:                     nil,
		MaxConnsPerHost:               20,
		MaxIdleConnDuration:           time.Minute,
		MaxConnDuration:               time.Minute,
		MaxIdemponentCallAttempts:     3,
		ReadBufferSize:                0,
		WriteBufferSize:               0,
		ReadTimeout:                   time.Second * 5,
		WriteTimeout:                  time.Second * 5,
		MaxResponseBodySize:           0,
		DisableHeaderNamesNormalizing: false,
		DisablePathNormalizing:        true,
		MaxConnWaitTimeout:            time.Second * 10,
		RetryIf:                       nil,
		ConnPoolStrategy:              fasthttp.LIFO,
		ConfigureClient:               nil,
	}
}

func call(method, uri string, timeout time.Duration, req, rsp interface{}) error {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	
	request.SetRequestURI(uri)
	request.Header.SetMethod(method)
	
	request.Header.SetContentType("application/json")
	
	request.Header.Set("X-Honoka-Rpc-Content-Type", "application/json")
	buf, err := sonic.Marshal(req)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	
	request.SetBodyRaw(buf)
	
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	
	err = client.DoTimeout(request, resp, timeout)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	
	switch resp.StatusCode() {
	case http.StatusOK:
		switch x := rsp.(type) {
		case *ext.ExtRsp:
			x.Buf = resp.Body()
		default:
			var baseRsp BaseRsp
			err = sonic.Unmarshal(resp.Body(), &baseRsp)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			
			if baseRsp.Code != 0 {
				return CreateErrorWithMsg(baseRsp.Code, baseRsp.Message)
			}
			
			err = sonic.Unmarshal([]byte(utils.ToString(baseRsp.Data)), rsp)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		}
	default:
		return CreateErrorWithMsg(-1, "http status code %d", resp.StatusCode())
	}
	
	return nil
}

func CallWithAddr(serviceName string, addr string, path string, req, rsp interface{}) error {
	return call(fasthttp.MethodPost, fmt.Sprintf("%s/%s", strings.TrimSuffix(addr, "/"), strings.TrimPrefix(path, "/")), time.Second*5, req, rsp)
}

func CallWithAddrMethod(serviceName string, method Method, addr string, path string, req, rsp interface{}) error {
	var m string
	switch method {
	case GET:
		m = fasthttp.MethodGet
	case HEAD:
		m = fasthttp.MethodHead
	case POST:
		m = fasthttp.MethodPost
	case PUT:
		m = fasthttp.MethodPut
	case PATCH:
		m = fasthttp.MethodPatch
	case DELETE:
		m = fasthttp.MethodDelete
	case CONNECT:
		m = fasthttp.MethodConnect
	case OPTIONS:
		m = fasthttp.MethodOptions
	case TRACE:
		m = fasthttp.MethodTrace
	default:
		return errors.New("unsupported method")
	}
	
	return call(m, fmt.Sprintf("%s/%s", strings.TrimSuffix(addr, "/"), strings.TrimPrefix(path, "/")), time.Second*5, req, rsp)
}
