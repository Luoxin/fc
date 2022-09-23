package parse

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/elliotchance/pie/v2"
	"github.com/emicklei/proto"
)

const (
	BaseModPath = "github.com/Luoxin/sexy"
	BaseDomain  = "github.com"
)

var protoMap = map[string]*Proto{}

type (
	Proto struct {
		MessageMap     map[string]*ProtoMessage
		EnumMap        map[string]*ProtoEnum
		Package        *ProtoPackage
		Service        *ProtoService
		GoPackagePath  string
		PackageDirPath string
	}
	
	ProtoMessage struct {
		source *proto.Message
		
		Name string
		
		FieldMap map[string]*ProtoFiled
		Comment  *ProtoCommentTag
	}
	
	ProtoFiled struct {
		source proto.Visitee
		
		Name    string
		typ     ProtoFiledType
		Comment *ProtoCommentTag
	}
	
	ProtoFiledType uint8
	
	ProtoCommentTag struct {
		comment *proto.Comment
		TagMap  map[string][]string
	}
	
	ProtoEnum struct {
		source   *proto.Enum
		Name     string
		Comment  *ProtoCommentTag
		FieldMap map[string]*ProtoEnumField
	}
	
	ProtoEnumField struct {
		source *proto.EnumField
		Name   string
		Value  int
		
		// 写在上面的注释
		Comment *ProtoCommentTag
		
		// 写在同行的注释
		InlineComment *ProtoCommentTag
	}
	
	ProtoPackage struct {
		source *proto.Package
		
		Name          string
		Comment       *ProtoCommentTag
		InlineComment *ProtoCommentTag
	}
	
	ProtoService struct {
		source  *proto.Service
		Name    string
		Comment *ProtoCommentTag
		RpcMap  map[string]*ProtoRPC
	}
	
	ProtoRPC struct {
		source       *proto.RPC
		Comment      *ProtoCommentTag
		RequestType  string
		ResponseType string
	}
)

func ParseProto(name string) (*Proto, error) {
	p := protoMap[name]
	if p != nil {
		return p, nil
	}
	
	p = &Proto{
		MessageMap: map[string]*ProtoMessage{},
		EnumMap:    map[string]*ProtoEnum{},
	}
	
	protoBuf, err := utils.FileRead(filepath.Join(ProtoDir(), name+".proto"))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	parser := proto.NewParser(bytes.NewBuffer(protoBuf))
	
	pb, err := parser.Parse()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	
	proto.Walk(
		pb,
		proto.WithMessage(p.AddMessage),
		proto.WithEnum(p.AddEnum),
		proto.WithPackage(p.AddPackage),
		proto.WithService(p.AddService),
		proto.WithOption(func(option *proto.Option) {
			switch option.Name {
			case "go_package":
				p.GoPackagePath = option.Constant.Source
				p.PackageDirPath = strings.TrimPrefix(p.GoPackagePath, BaseModPath)
			}
		}),
	)
	
	protoMap[name] = p
	return p, nil
}

func (p *Proto) AddMessage(message *proto.Message) {
	m := NewProtoMessage(message)
	p.MessageMap[m.Name] = m
}

func (p *Proto) AddEnum(enum *proto.Enum) {
	m := NewProtoEnum(enum)
	p.EnumMap[m.Name] = m
}

func (p *Proto) AddPackage(pkg *proto.Package) {
	p.Package = NewProtoPackage(pkg)
}

func (p *Proto) AddService(service *proto.Service) {
	p.Service = NewProtoService(service)
}

func NewProtoMessage(message *proto.Message) *ProtoMessage {
	p := &ProtoMessage{
		source:   message,
		Name:     message.Name,
		FieldMap: map[string]*ProtoFiled{},
		Comment:  NewCommentTag(message.Comment),
	}
	
	parent := message.Parent
	if parent != nil {
		switch x := parent.(type) {
		case *proto.Proto:
			// 最上层了
		case *proto.Message:
			m := NewProtoMessage(x)
			p.Name = fmt.Sprintf("%s_%s", m.Name, p.Name)
		default:
			log.Warnf("unknown %v", reflect.TypeOf(message.Parent).Elem())
		}
	}
	
	for _, element := range message.Elements {
		filed := NewProtoFiled(element)
		p.FieldMap[filed.Name] = filed
	}
	
	return p
}

func NewCommentTag(comment *proto.Comment) *ProtoCommentTag {
	p := &ProtoCommentTag{
		comment: comment,
		TagMap:  map[string][]string{},
	}
	
	if comment != nil {
		for _, line := range comment.Lines {
			line := strings.TrimSpace(line)
			if !strings.HasPrefix(line, "@") {
				continue
			}
			
			line = strings.TrimPrefix(line, "@")
			
			idx := strings.Index(line, ":")
			if idx < 0 {
				continue
			}
			
			tag := line[:idx]
			if _, ok := p.TagMap[tag]; !ok {
				p.TagMap[tag] = []string{}
			}
			
			p.TagMap[tag] = append(p.TagMap[tag], line[idx+1:])
		}
	}
	
	return p
}

func (p *ProtoCommentTag) Exist(key string) bool {
	_, ok := p.TagMap[key]
	return ok
}

func (p *ProtoCommentTag) Get(key string) map[string]string {
	m := map[string]string{}
	pie.Each[string](strings.Split(p.GetString(key), ";"), func(item string) {
		idx := strings.Index(item, "=")
		if idx < 0 {
			m[item] = ""
			return
		}
		
		m[item[:idx]] = item[idx+1:]
	})
	return m
}

func (p *ProtoCommentTag) GetString(key string) string {
	return strings.Join(p.TagMap[key], ";")
}

const (
	ProtoFiledTypeMessage ProtoFiledType = iota + 1
	ProtoFiledTypeNormalField
	ProtoFiledTypeEnumField
	ProtoFiledTypeEnum
	ProtoFiledTypeMapField
	ProtoFiledTypeOneOfField
)

func NewProtoFiled(filed proto.Visitee) *ProtoFiled {
	p := &ProtoFiled{
		source: filed,
	}
	
	switch x := filed.(type) {
	case *proto.Message:
		// m := NewProtoMessage(x)
		p.Name = x.Name
		p.typ = ProtoFiledTypeMessage
		p.Comment = NewCommentTag(x.Comment)
	case *proto.NormalField:
		p.Name = x.Name
		p.typ = ProtoFiledTypeNormalField
		p.Comment = NewCommentTag(x.Comment)
	case *proto.EnumField:
		p.Name = x.Name
		p.typ = ProtoFiledTypeEnumField
		p.Comment = NewCommentTag(x.Comment)
	case *proto.Enum:
		p.Name = x.Name
		p.typ = ProtoFiledTypeEnum
		p.Comment = NewCommentTag(x.Comment)
	case *proto.MapField:
		p.Name = x.Name
		p.typ = ProtoFiledTypeMapField
		p.Comment = NewCommentTag(x.Comment)
	case *proto.OneOfField:
		p.Name = x.Name
		p.typ = ProtoFiledTypeOneOfField
		p.Comment = NewCommentTag(x.Comment)
	default:
		log.Warnf("unknown %v", reflect.TypeOf(filed).Elem())
	}
	
	return p
}

func NewProtoEnum(enum *proto.Enum) *ProtoEnum {
	p := &ProtoEnum{
		source:   enum,
		Name:     enum.Name,
		Comment:  NewCommentTag(enum.Comment),
		FieldMap: map[string]*ProtoEnumField{},
	}
	
	parent := enum.Parent
	switch x := parent.(type) {
	case *proto.Proto:
		// 是最上层了
	case *proto.Message:
		m := NewProtoMessage(x)
		p.Name = fmt.Sprintf("%s_%s", m.Name, p.Name)
	default:
		log.Warnf("unknown %v", reflect.TypeOf(parent).Elem())
	}
	
	for _, element := range enum.Elements {
		switch x := element.(type) {
		case *proto.EnumField:
			field := NewProtoEnumField(x)
			p.FieldMap[field.Name] = field
		default:
			log.Warnf("unknown %v", reflect.TypeOf(element).Elem())
		}
	}
	
	return p
}

func NewProtoEnumField(field *proto.EnumField) *ProtoEnumField {
	p := &ProtoEnumField{
		source:        field,
		Name:          field.Name,
		Value:         field.Integer,
		Comment:       NewCommentTag(field.Comment),
		InlineComment: NewCommentTag(field.InlineComment),
	}
	
	return p
}

func NewProtoPackage(pkg *proto.Package) *ProtoPackage {
	p := &ProtoPackage{
		source:        pkg,
		Name:          pkg.Name,
		Comment:       NewCommentTag(pkg.Comment),
		InlineComment: NewCommentTag(pkg.InlineComment),
	}
	
	return p
}

func NewProtoService(service *proto.Service) *ProtoService {
	p := &ProtoService{
		source:  service,
		Name:    service.Name,
		Comment: NewCommentTag(service.Comment),
		RpcMap:  map[string]*ProtoRPC{},
	}
	
	for _, element := range service.Elements {
		switch x := element.(type) {
		case *proto.RPC:
			p.RpcMap[x.Name] = NewProtoRPC(x)
		default:
			log.Warnf("unknown %v", reflect.TypeOf(element).Elem())
		}
	}
	
	return p
}

func NewProtoRPC(rpc *proto.RPC) *ProtoRPC {
	p := &ProtoRPC{
		source:       rpc,
		Comment:      NewCommentTag(rpc.Comment),
		RequestType:  rpc.RequestType,
		ResponseType: rpc.ReturnsType,
	}
	
	return p
}
