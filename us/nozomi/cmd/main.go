package main

import (
	"github.com/Luoxin/sexy/honoka"
	"github.com/Luoxin/sexy/us/nozomi"
	"github.com/Luoxin/sexy/us/nozomi/impl"
	"github.com/Luoxin/sexy/us/rpc"
	"github.com/darabuchi/log"
)

func init() {
	honoka.ServiceName = nozomi.ServiceName
}

func main() {
	err := honoka.Load()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	
	rpc.RegisterHandlers(apiList...)
	
	err = impl.InitState()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	
	err = rpc.StartService()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
