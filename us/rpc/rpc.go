package rpc

import (
	"github.com/valyala/fasthttp"
)

type Ctx struct {
	source *fasthttp.RequestCtx
}

func (p *Ctx) GetPath() string {
	if p == nil {
		return ""
	}
	
	if p.source == nil {
		return ""
	}
	
	return string(p.source.Path())
}

func (p *Ctx) Method() Method {
	if p == nil {
		return MethodNone
	}
	
	if p.source == nil {
		return MethodNone
	}
	
	switch string(p.source.Method()) {
	case fasthttp.MethodGet:
		return GET
	case fasthttp.MethodHead:
		return HEAD
	case fasthttp.MethodPost:
		return POST
	case fasthttp.MethodPut:
		return PUT
	case fasthttp.MethodPatch:
		return PATCH
	case fasthttp.MethodDelete:
		return CONNECT
	case fasthttp.MethodConnect:
		return CONNECT
	case fasthttp.MethodOptions:
		return OPTIONS
	case fasthttp.MethodTrace:
		return TRACE
	}
	
	return MethodNone
}

type (
	BaseRsp struct {
		Error
		Hint string `json:"hint,omitempty"`
		Data any    `json:"data,omitempty"`
	}
)
