package rpc

import (
	"bytes"
	"reflect"
	"strconv"

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

	if ft.NumIn() != 2 {
		log.Panicf("num in not 2")
	}

	// 第0位为ctx
	x := ft.In(0)
	for x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	if x.Name() != "Ctx" {
		log.Panicf("first in not ctx")
	}

	// 第二位要是struct
	{
		x := ft.In(1)
		for x.Kind() == reflect.Ptr {
			x = x.Elem()
		}
		if x.Kind() != reflect.Struct {
			log.Panicf("second in not struct")
		}

		h.req = x
	}

	if ft.NumOut() != 2 {
		log.Panicf("num out not 2")
	}

	{
		x := ft.Out(0)
		for x.Kind() == reflect.Ptr {
			x = x.Elem()
		}
		if x.Kind() != reflect.Struct {
			log.Panicf("first out not struct")
		}

		h.rsp = x
	}

	x = ft.Out(1)
	if x.Name() != "error" {
		log.Panicf("second out not error")
	}

	return h
}

func RegisterHandler(method Method, path string, handler any) {
	h := NewHandler(handler)

	logic := func(ctx *fasthttp.RequestCtx) {
		log.SetTrace(log.GenTraceId())
		defer log.DelTrace()

		utils.CachePanicWithHandle(func(err any) {

		})
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

		resp, err := func() (any, *Error) {
			obj := reflect.New(h.req)
			if _, ok := obj.Interface().(*ExtReq); !ok {
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

			ret := h.f.Call([]reflect.Value{
				reflect.ValueOf(&Ctx{}),
				obj,
			})
			if len(ret) != 2 {
				log.Panicf("num out not 2")
			}

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
		}()
		if err != nil {
			_, err := ctx.Write(err.Marshal())
			if err != nil {
				log.Errorf("err:%v", err)
			}
		} else {
			if ext, ok := resp.(*ExtRsp); ok {
				_, err := ctx.Write(ext.Buf)
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

	switch method {
	case GET:
		handlerRouter.Handle(fasthttp.MethodGet, path, logic)
	case HEAD:
		handlerRouter.Handle(fasthttp.MethodHead, path, logic)
	case POST:
		handlerRouter.Handle(fasthttp.MethodPost, path, logic)
	case PUT:
		handlerRouter.Handle(fasthttp.MethodPut, path, logic)
	case PATCH:
		handlerRouter.Handle(fasthttp.MethodPatch, path, logic)
	case DELETE:
		handlerRouter.Handle(fasthttp.MethodDelete, path, logic)
	case CONNECT:
		handlerRouter.Handle(fasthttp.MethodConnect, path, logic)
	case OPTIONS:
		handlerRouter.Handle(fasthttp.MethodOptions, path, logic)
	case TRACE:
		handlerRouter.Handle(fasthttp.MethodTrace, path, logic)
	case All:
		handlerRouter.Handle(fasthttp.MethodGet, path, logic)
		handlerRouter.Handle(fasthttp.MethodHead, path, logic)
		handlerRouter.Handle(fasthttp.MethodPost, path, logic)
		handlerRouter.Handle(fasthttp.MethodPut, path, logic)
		handlerRouter.Handle(fasthttp.MethodPatch, path, logic)
		handlerRouter.Handle(fasthttp.MethodDelete, path, logic)
		handlerRouter.Handle(fasthttp.MethodConnect, path, logic)
		handlerRouter.Handle(fasthttp.MethodOptions, path, logic)
		handlerRouter.Handle(fasthttp.MethodTrace, path, logic)
	}
}
