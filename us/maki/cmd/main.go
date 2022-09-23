package main

import (
	"io/fs"
	"path/filepath"
	"strings"
	
	"github.com/Luoxin/sexy/us/maki/gen"
	"github.com/Luoxin/sexy/us/maki/parse"
	"github.com/darabuchi/log"
)

func main() {
	err := filepath.Walk(parse.ProtoDir(), func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) != ".proto" {
			return nil
		}
		
		name := strings.TrimSuffix(filepath.Base(path), ".proto")
		
		_, err = parse.ParseProto(name)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		
		err = gen.GenPb(name)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		
		return nil
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	
	parse.ParseGoFile("E:\\go\\darabuchi\\sexy\\us\\maki\\parse\\proto.go")
}
