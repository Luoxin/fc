package gen

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	
	"github.com/Luoxin/sexy/us/maki/parse"
	"github.com/darabuchi/log"
)

func GenPb(name string) error {
	log.Infof("try gen %s", name)
	
	proto, err := parse.ParseProto(name)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	
	cmd := exec.Command(
		filepath.Join(parse.Root(), "bin", ProtocName),
		"--plugin=protoc-gen-go="+filepath.Join(parse.Root(), "bin", ProtocGenGoName),
		"--proto_path="+filepath.Join(parse.Root(), "proto"),
		"--go_out="+parse.Root(),
		filepath.Join(parse.Root(), "proto", name+".proto"),
	)
	cmd.Dir = parse.Root()
	cmd.Env = os.Environ()
	
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	
	err = cmd.Run()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	
	defer func() {
		err = os.RemoveAll(filepath.Join(parse.Root(), parse.BaseDomain))
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	}()
	
	source, err := os.Open(filepath.Join(parse.Root(), proto.GoPackagePath, name+".pb.go"))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer source.Close()
	
	target, err := os.OpenFile(filepath.Join(parse.Root(), proto.PackageDirPath, name+".pb.go"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer target.Close()
	
	_, err = io.Copy(target, source)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	
	return nil
}
