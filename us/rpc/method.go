package rpc

type Method uint32

const (
	MethodNone        = 0 << 1
	GET        Method = 1 << 1
	POST       Method = 1 << 2
	PUT        Method = 1 << 3
	PATCH      Method = 1 << 4
	DELETE     Method = 1 << 5
	// COPY      Method = 1 << 6
	HEAD    Method = 1 << 7
	OPTIONS Method = 1 << 8
	// LINK      Method = 1 << 9
	// UNLINK    Method = 1 << 10
	// PURGE     Method = 1 << 11
	// LOCK      Method = 1 << 12
	// UNLOCK    Method = 1 << 13
	// PROPFIND  Method = 1 << 14
	// VIEW      Method = 1 << 15
	CONNECT Method = 1 << 16
	TRACE   Method = 1 << 17
	
	All Method = 0xffff
)
