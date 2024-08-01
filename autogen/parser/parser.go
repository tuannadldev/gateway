package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gateway/autogen/model"

	"github.com/emicklei/proto"
)

func ParseProtoFile(path string) *model.Proto {
	reader, _ := os.Open(path)
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	visitor := new(protoVisitor)
	definition.Accept(visitor)

	proto := new(model.Proto)
	proto.FilePath = path
	proto.Package = visitor.packageName
	proto.Services = make(map[string]*model.Service)
	proto.Messages = make(map[string]*model.Message)
	for _, service := range visitor.services {
		proto.Services[service.Name] = service
	}
	for _, message := range visitor.messages {
		proto.Messages[message.Name] = message
	}

	return proto
}

func MergeProtoMessages(protos ...*model.Proto) map[string]*model.Message {
	messages := make(map[string]*model.Message)
	for _, proto := range protos {
		for name, message := range proto.Messages {
			if _, hasKey := messages[name]; hasKey {
				fmt.Println("Found duplicated message with name:", name, ", skip message definition in package:", proto.Package)
				continue
			}
			messages[name] = message
		}
	}

	return messages
}

type protoVisitor struct {
	proto.NoopVisitor
	packageName string
	services    []*model.Service
	messages    []*model.Message
}

func (pv *protoVisitor) VisitPackage(p *proto.Package) {
	pv.packageName = p.Name
}

func (pv *protoVisitor) VisitService(s *proto.Service) {
	if s.Comment == nil {
		fmt.Println("Comment not found, skip service:", s.Name)
		return
	}

	av := new(apiVisitor)
	for _, element := range s.Elements {
		element.Accept(av)
	}
	service := model.Service{
		Name:    s.Name,
		Path:    strings.TrimSpace(s.Comment.Lines[0]),
		Package: pv.packageName,
		Apis:    av.apis,
	}
	pv.services = append(pv.services, &service)
}

func (pv *protoVisitor) VisitMessage(m *proto.Message) {
	fv := new(fieldVisitor)
	for _, element := range m.Elements {
		element.Accept(fv)
	}
	message := model.Message{
		Name:   m.Name,
		Fields: fv.fields,
	}
	pv.messages = append(pv.messages, &message)
}

type fieldVisitor struct {
	proto.NoopVisitor
	fields []*model.Field
}

func (fv *fieldVisitor) VisitNormalField(f *proto.NormalField) {
	desc := ""
	if f.Comment != nil {
		desc = strings.Join(f.Comment.Lines[:], " ")
	} else {
		desc = "Missing comment! you should fill its for long-term"
	}
	field := model.Field{
		Name:        f.Name,
		Type:        f.Type,
		Bit:         parseIntBit(f.Type),
		Location:    "query",
		Required:    false,
		Repeated:    f.Repeated,
		Description: desc,
	}
	if f.InlineComment != nil {
		field.Required = strings.Contains(f.InlineComment.Message(), "required")
		if strings.Contains(f.InlineComment.Message(), "path") {
			field.Location = "path"
		}
	}
	fv.fields = append(fv.fields, &field)
}

type apiVisitor struct {
	proto.NoopVisitor
	apis []*model.Api
}

func (av *apiVisitor) VisitRPC(r *proto.RPC) {
	if r.Comment == nil {
		fmt.Println("Comment not found, skip api:", r.Name)
		return
	}

	tokens := strings.Split(r.Comment.Lines[0], ",")
	api := model.Api{
		Name:         r.Name,
		Method:       strings.ToUpper(strings.TrimSpace(tokens[0])),
		Path:         strings.ToLower(strings.TrimSpace(tokens[1])),
		RequestType:  r.RequestType,
		ResponseType: r.ReturnsType,
	}
	if len(tokens) > 2 {
		// api.Permission = strings.ToLower(strings.TrimSpace(tokens[2]))
		api.Permission = strings.TrimSpace(tokens[2])
	}
	av.apis = append(av.apis, &api)
}

func parseIntBit(typeName string) int64 {
	if strings.HasPrefix(typeName, "int") {
		if len(typeName[3:]) > 0 {
			bit, err := strconv.ParseInt(typeName[3:], 10, 32)
			if err != nil {
				fmt.Println("failed to parse int bit of type:", typeName)
				return 0
			}
			return bit
		}
		return 32
	}
	return 0
}
