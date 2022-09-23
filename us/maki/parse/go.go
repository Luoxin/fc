package parse

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
)

type (
	GoFile struct {
		FilePath    string
		PackageName string
		source      *ast.File
		ImportMap   map[string]*GoImport
		FuncMap     map[string]*GoFunc
	}
	
	GoImport struct {
		source *ast.ImportSpec
		
		Name string
		Path string
	}
	
	GoFunc struct {
		Name   string
		source *ast.FuncDecl
	}
)

var (
	ErrFileNotFound = errors.New("file not found")
)

func ParseGoDir(path string) error {
	fSet := token.NewFileSet()
	m, err := parser.ParseDir(fSet, path, nil, 0)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	
	for k, v := range m {
		log.Info(k)
		log.Info(v.Files)
		log.Info(v.Imports)
	}
	
	return nil
}

func ParseGoFile(path string) (*GoFile, error) {
	if !utils.IsFile(path) {
		log.Warnf("%s is not found", path)
		return nil, ErrFileNotFound
	}
	
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, path, nil, 0)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	
	gf := &GoFile{
		source:      f,
		FilePath:    path,
		PackageName: f.Name.Name,
		ImportMap:   map[string]*GoImport{},
		FuncMap:     map[string]*GoFunc{},
	}
	
	for _, spec := range f.Imports {
		i := NewGoImport(spec)
		gf.ImportMap[i.Path] = i
	}
	
	for _, decl := range f.Decls {
		switch x := decl.(type) {
		case *ast.FuncDecl:
			gf.FuncMap[x.Name.Name] = NewGoFunc(x)
		case *ast.GenDecl:
		
		default:
			log.Warnf("unknown %v", reflect.TypeOf(decl).Elem())
		}
	}
	
	return gf, nil
}

func NewGoImport(source *ast.ImportSpec) *GoImport {
	p := &GoImport{
		source: source,
	}
	
	if source.Name != nil {
		p.Name = source.Name.Name
	}
	
	if source.Path != nil {
		p.Path = source.Path.Value
	}
	
	return p
}

func NewGoFunc(source *ast.FuncDecl) *GoFunc {
	p := &GoFunc{
		source: source,
		Name:   source.Name.Name,
	}
	
	if source.Recv != nil {
		log.Info(source.Recv.List[0].Names)
		log.Info(source.Recv.List[0].Comment)
	}
	
	return p
}
