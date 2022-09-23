package rpc

import (
	"bytes"
	"reflect"
	"strconv"
	
	"github.com/Luoxin/sexy/us/ext"
	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

var handlerRouter = router.New()

type Handler struct {
	f        reflect.Value
	req, rsp reflect.Type
	
	reqCnt, rspCnt int
}

func NewHandler(handler any) *Handler {
	h := &Handler{
		f: reflect.ValueOf(handler),
	}
	
	f := h.f
	ft := f.Type()
	
	if ft.Kind() != reflect.Func {
		log.Panicf("not func")
	}
	
	h.parseReq()
	h.parseRsp()
	
	return h
}

func (p *Handler) parseReq() {
	f := p.f
	ft := f.Type()
	
	p.reqCnt = ft.NumIn()
	
	// 第0位为ctx
	x := ft.In(0)
	for x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	
	if x.Name() != "Ctx" {
		log.Panicf("first in not ctx")
	}
	
	if p.reqCnt > 1 {
		// 第二位如果存在 必须为struct
		x := ft.In(1)
		for x.Kind() == reflect.Ptr {
			x = x.Elem()
		}
		if x.Kind() != reflect.Struct {
			log.Panicf("second in not struct")
		}
		
		p.req = x
	}
	
	if p.reqCnt > 2 {
		panic("unsupported return num")
	}
}

func (p *Handler) parseRsp() {
	f := p.f
	ft := f.Type()
	
	p.rspCnt = ft.NumOut()
	
	switch p.rspCnt {
	case 1: // 如果就一个返回值，那么他一定要是 error 类型
		x := ft.Out(0)
		if x.Name() != "error" {
			log.Panicf("second out not error")
		}
	case 2:
		{
			x := ft.Out(0)
			for x.Kind() == reflect.Ptr {
				x = x.Elem()
			}
			if x.Kind() != reflect.Struct {
				log.Panicf("first out not struct")
			}
			
			p.rsp = x
		}
		
		x := ft.Out(1)
		if x.Name() != "error" {
			log.Panicf("second out not error")
		}
	default:
		panic("unsupported return num")
	}
}

func (h *Handler) HandleFastHttp(ctx *fasthttp.RequestCtx) {
	log.SetTrace(log.GenTraceId())
	defer log.DelTrace()
	
	utils.CachePanicWithHandle(func(err any) {
	
	})
	
	// 日志部分
	b := bytes.NewBuffer(nil)
	
	b.Write(ctx.Method())
	b.WriteString(" ")
	b.Write(ctx.Path())
	b.WriteString(" ")
	
	b.WriteString(ctx.LocalAddr().String())
	b.WriteString("->")
	b.WriteString(ctx.RemoteAddr().String())
	b.WriteString(" ")
	
	if ctx.QueryArgs() != nil {
		b.Write(ctx.QueryArgs().QueryString())
		b.WriteString(" ")
	}
	
	if ctx.PostArgs() != nil {
		b.Write(ctx.PostArgs().QueryString())
		b.WriteString(" ")
	}
	
	if len(ctx.PostBody()) > 0 {
		b.WriteString(utils.ShortStr4Web(strconv.Quote(string(ctx.PostBody())), 1024*1024*10))
		b.WriteString(" ")
	}
	
	log.Info(b.String())
	defer func() {
		if len(ctx.Response.Body()) > 0 {
			b.WriteString("resp:")
			b.WriteString(utils.ShortStr4Web(strconv.Quote(string(ctx.Response.Body())), 1024*1024))
			b.WriteString(" ")
		}
		
		b.WriteString("status:")
		b.WriteString(strconv.Itoa(ctx.Response.StatusCode()))
		b.WriteString(" ")
		
		log.Info(b.String())
	}()
	
	// 请求部分
	resp, err := func() (any, *Error) {
		var callReq []reflect.Value
		
		callReq = append(callReq, reflect.ValueOf(&Ctx{
			source: ctx,
		}))
		
		switch h.reqCnt {
		case 2:
			obj := reflect.New(h.req)
			if _, ok := obj.Interface().(*ext.ExtReq); !ok {
				if len(ctx.PostBody()) > 0 {
					err := sonic.Unmarshal(ctx.PostBody(), obj.Interface())
					if err != nil {
						log.Errorf("err:%v", err)
						return nil, &Error{
							Code:    -1,
							Message: err.Error(),
						}
					}
				}
			}
			callReq = append(callReq, obj)
		}
		
		ret := h.f.Call(callReq)
		
		switch h.rspCnt {
		case 1:
			inter := ret[0].Interface()
			if inter != nil {
				log.Errorf("err:%v", inter)
				
				if x, ok := inter.(*Error); ok {
					return nil, x
				}
				
				return nil, &Error{
					Code:    -1,
					Message: inter.(error).Error(),
				}
			}
			
			return nil, nil
		case 2:
			inter := ret[1].Interface()
			
			if inter != nil {
				log.Errorf("err:%v", inter)
				
				if x, ok := inter.(*Error); ok {
					return nil, x
				}
				
				return nil, &Error{
					Code:    -1,
					Message: inter.(error).Error(),
				}
			}
			return ret[0].Interface(), nil
		default:
			return nil, &Error{
				Code:    -1,
				Message: "unsupported return cnt",
			}
		}
	}()
	if err != nil {
		log.Errorf("err:%v", err)
		_, err := ctx.Write(err.Marshal())
		if err != nil {
			log.Errorf("err:%v", err)
		}
	} else if resp != nil {
		if extRsp, ok := resp.(*ext.ExtRsp); ok {
			_, err := ctx.Write(extRsp.Buf)
			if err != nil {
				log.Errorf("err:%v", err)
			}
		} else {
			buf, err := sonic.Marshal(&BaseRsp{
				Error: Error{},
				Data:  resp,
			})
			if err != nil {
				log.Errorf("err:%v", err)
			} else {
				_, err := ctx.Write(buf)
				if err != nil {
					log.Errorf("err:%v", err)
				}
			}
		}
	}
}

type ApiHandler struct {
	Method  Method
	Path    string
	Handler any
}

func RegisterHandlers(apis ...ApiHandler) {
	for _, api := range apis {
		RegisterHandler(api)
	}
}

func RegisterHandler(api ApiHandler) {
	h := NewHandler(api.Handler)
	
	switch api.Method {
	case GET:
		handlerRouter.Handle(fasthttp.MethodGet, api.Path, h.HandleFastHttp)
	case HEAD:
		handlerRouter.Handle(fasthttp.MethodHead, api.Path, h.HandleFastHttp)
	case POST:
		handlerRouter.Handle(fasthttp.MethodPost, api.Path, h.HandleFastHttp)
	case PUT:
		handlerRouter.Handle(fasthttp.MethodPut, api.Path, h.HandleFastHttp)
	case PATCH:
		handlerRouter.Handle(fasthttp.MethodPatch, api.Path, h.HandleFastHttp)
	case DELETE:
		handlerRouter.Handle(fasthttp.MethodDelete, api.Path, h.HandleFastHttp)
	case CONNECT:
		handlerRouter.Handle(fasthttp.MethodConnect, api.Path, h.HandleFastHttp)
	case OPTIONS:
		handlerRouter.Handle(fasthttp.MethodOptions, api.Path, h.HandleFastHttp)
	case TRACE:
		handlerRouter.Handle(fasthttp.MethodTrace, api.Path, h.HandleFastHttp)
	case All:
		handlerRouter.Handle(fasthttp.MethodGet, api.Path, h.HandleFastHttp)
		handlerRouter.Handle(fasthttp.MethodHead, api.Path, h.HandleFastHttp)
		handlerRouter.Handle(fasthttp.MethodPost, api.Path, h.HandleFastHttp)
		handlerRouter.Handle(fasthttp.MethodPut, api.Path, h.HandleFastHttp)
		handlerRouter.Handle(fasthttp.MethodPatch, api.Path, h.HandleFastHttp)
		handlerRouter.Handle(fasthttp.MethodDelete, api.Path, h.HandleFastHttp)
		handlerRouter.Handle(fasthttp.MethodConnect, api.Path, h.HandleFastHttp)
		handlerRouter.Handle(fasthttp.MethodOptions, api.Path, h.HandleFastHttp)
		handlerRouter.Handle(fasthttp.MethodTrace, api.Path, h.HandleFastHttp)
	}
}
