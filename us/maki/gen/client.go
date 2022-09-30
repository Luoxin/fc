package gen

import (
	"path/filepath"
	
	"github.com/Luoxin/sexy/us/maki/parse"
	"github.com/darabuchi/log"
)

func GenClient(name string) error {
	pb, err := parse.ParseProto(name)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	
	if pb.Service == nil {
		log.Warnf("not found service, skip gen client")
		return nil
	}
	
	parse.ParseGoFile(filepath.Join(parse.Root(), pb.PackageDirPath, name+".pb.go"))
	
	return nil
}
