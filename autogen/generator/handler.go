package generator

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"gateway/application/routing"
	"gateway/autogen/model"
	"gateway/autogen/parser"
	"gateway/autogen/util"
)

const HANDLER_PARENT_PATH = "application/routing/delivery/handler"
const ROUTING_REGISTRY_PATH = "application/routing/delivery/routingRegistry.go"

func GenerateHandlers(protos []*model.Proto) ([]*model.Handler, []string, error) {
	// because messages can be defined in a separate proto file
	// we need to create a map of all messages and use it for generating handler for each service
	messages := parser.MergeProtoMessages(protos...)
	var allHandlers []*model.Handler
	handlerFolderPaths := make(map[string]bool) // use a map to collect non-duplicated handler folder paths
	for _, proto := range protos {
		handlers, err := generateHandlers(proto, messages)
		if err != nil {
			fmt.Println("failed to generate handler for:", proto.FilePath, err)
		} else {
			allHandlers = append(allHandlers, handlers...)
			for _, handler := range handlers {
				handlerFolderPaths[handler.FolderPath] = true
			}
		}
	}

	folderPaths := make([]string, 0, len(handlerFolderPaths))
	for path := range handlerFolderPaths {
		folderPaths = append(folderPaths, path)
	}

	return allHandlers, folderPaths, nil
}

func generateHandlers(proto *model.Proto, messages map[string]*model.Message) ([]*model.Handler, error) {
	// proto.FilePath must follow this format: proto/something/something.proto
	var segments []string
	if runtime.GOOS == "windows" {
		segments = strings.Split(proto.FilePath, "\\")
	} else {
		segments = strings.Split(proto.FilePath, "/")
	}
	if len(segments) < 3 {
		return nil, fmt.Errorf("unsupported proto file path: %s", proto.FilePath)
	}
	parentFolderPath := HANDLER_PARENT_PATH
	for i := 1; i < len(segments)-1; i++ {
		folderPath := parentFolderPath + "/" + segments[i]
		err := util.CreateFolderIfNotExist(folderPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create new folder at: %s, skip generating handler, error: %e", folderPath, err)
		}
		parentFolderPath = folderPath
	}
	if segments[len(segments)-2] != proto.Package {
		parentFolderPath += "/" + proto.Package
		err := util.CreateFolderIfNotExist(parentFolderPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create new folder at: %s, skip generating handler, error: %e", parentFolderPath, err)
		}
	}

	handlers := []*model.Handler{}
	generationErrors := []error{}
	for _, service := range proto.Services {
		fmt.Println(service)

		for _, api := range service.Apis {
			fmt.Println("Gen code for handler:", api.Name)
			handler := model.Handler{
				Name:         api.Name,
				Method:       api.Method,
				Path:         service.Path + api.Path,
				Permission:   api.Permission,
				ServiceName:  service.Name,
				ProtoFolder:  strings.Join(segments[1:len(segments)-1], "/"),
				ProtoPackage: service.Package,
				RequestType:  api.RequestType,
				ResponseType: api.ResponseType,
				FolderPath:   parentFolderPath,
			}
			var err error
			switch api.Method {
			case "GET":
				err = generateGetHandler(&handler, messages[api.RequestType])
			case "POST":
				err = generatePostHandler(&handler, messages[api.RequestType])
			case "PUT":
				err = generatePutHandler(&handler, messages[api.RequestType])
			default:
				err = fmt.Errorf("failed to generate handler for: %s, unsupported method: %s", api.Name, api.Method)
			}
			if err != nil {
				generationErrors = append(generationErrors, err)
			} else {
				handlers = append(handlers, &handler)
			}
		}
	}

	var err error
	if len(generationErrors) > 0 {
		err = generationErrors[0]
		for _, e := range generationErrors[1:] {
			err = fmt.Errorf("%w; %w", err, e)
		}

		return nil, err
	}

	return handlers, nil
}

func generatePostHandler(handler *model.Handler, requestType *model.Message) error {
	template, err := util.ReadFileAsString("autogen/template/handler_post.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	template = strings.ReplaceAll(template, "<proto_folder>", handler.ProtoFolder)
	template = strings.ReplaceAll(template, "<proto_package>", handler.ProtoPackage)
	template = strings.ReplaceAll(template, "<handler_name>", handler.Name)
	handlerNameLower1st := util.LowerFirstChar(handler.Name)
	template = strings.ReplaceAll(template, "<handler_name_lower_1st>", handlerNameLower1st)
	template = strings.ReplaceAll(template, "<request_type>", handler.RequestType)
	template = strings.ReplaceAll(template, "<request_type_lower_1st>", util.LowerFirstChar(handler.RequestType))
	template = strings.ReplaceAll(template, "<response_type>", handler.ResponseType)
	template = strings.ReplaceAll(template, "<path>", handler.Path)
	if routing.ApiPublic[handler.Path] {
		template = strings.ReplaceAll(template, "<summary>", "api public")
	} else {
		template = strings.ReplaceAll(template, "<summary>", "permission: "+handler.Permission)
	}
	template = strings.ReplaceAll(template, "<tags>", handler.ServiceName)

	// handle params
	paramDecTemplate, err := util.ReadFileAsString("autogen/template/param_declaration.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingIdTemplate, err := util.ReadFileAsString("autogen/template/param_parsing_id.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramDeclarations := []string{}
	paramParsings := []string{}
	hasIntField := false
	for _, field := range requestType.Fields {

		declaration := paramDecTemplate
		declaration = strings.ReplaceAll(declaration, "<param_name>", field.Name)
		if field.Location == "path" {
			declaration = strings.ReplaceAll(declaration, "<param_location>", field.Location)
			declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(true))
		} else {
			declaration = strings.ReplaceAll(declaration, "<param_location>", "body")
			declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(field.Required))
		}

		field_type := field.Type
		if strings.Contains(field_type, "Request") {
			field_type = handler.ProtoPackage + "." + field_type
		} else {
			switch field_type {
			case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "complex64", "complex128", "byte", "rune", "string", "bool", "enum", "bytes":
			case "float", "double", "float64":
				field_type = "float32"
			default:
				field_type = handler.ProtoPackage + "." + field_type
			}
		}

		if field.Repeated {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "[]"+field_type)
		} else if field.Location == "enum" {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "string")
		} else if field.Location == "file" {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "file")
		} else {
			declaration = strings.ReplaceAll(declaration, "<param_type>", field_type)
		}
		declaration = strings.ReplaceAll(declaration, "<param_description>", field.Description)

		paramDeclarations = append(paramDeclarations, declaration)

		if field.Location == "path" {
			parsing := paramParsingIdTemplate
			hasIntField = true
			parsing = paramParsingIdTemplate

			parsing = strings.ReplaceAll(parsing, "<param_name>", field.Name)
			paramContext := "Query"
			if field.Location == "path" {
				paramContext = "Param"
			}
			parsing = strings.ReplaceAll(parsing, "<param_context>", paramContext)
			parsing = strings.ReplaceAll(parsing, "<param_bit>", strconv.FormatInt(field.Bit, 10))
			parsing = strings.ReplaceAll(parsing, "<param_name_upper_case_word>", util.ConvertSliceStrToUCWord(strings.Split(field.Name, "_")))
			paramParsings = append(paramParsings, parsing)
		}
	}
	if len(paramDeclarations) > 0 {
		template = strings.ReplaceAll(template, "<param_declarations>", strings.Join(paramDeclarations, "\n"))
	} else {
		template = strings.ReplaceAll(template, "\n<param_declarations>", strings.Join(paramDeclarations, "\n"))
	}
	template = strings.ReplaceAll(template, "<param_parsings>", strings.Join(paramParsings, "\n"))
	if !hasIntField {
		if runtime.GOOS == "windows" {
			regex := regexp.MustCompile(`"strconv"`)
			template = regex.ReplaceAllString(template, "")
		} else {
			regex := regexp.MustCompile(" +\"strconv\"\\n")
			template = regex.ReplaceAllString(template, "")
		}
	}

	err = os.WriteFile(handler.FolderPath+"/"+handlerNameLower1st+".go", []byte(template), 0644)
	if err != nil {
		fmt.Println("failed to write handler file:", err)
		return err
	}

	return nil
}

func generatePutHandler(handler *model.Handler, requestType *model.Message) error {
	template, err := util.ReadFileAsString("autogen/template/handler_put.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	template = strings.ReplaceAll(template, "<proto_folder>", handler.ProtoFolder)
	template = strings.ReplaceAll(template, "<proto_package>", handler.ProtoPackage)
	template = strings.ReplaceAll(template, "<handler_name>", handler.Name)
	handlerNameLower1st := util.LowerFirstChar(handler.Name)
	template = strings.ReplaceAll(template, "<handler_name_lower_1st>", handlerNameLower1st)
	template = strings.ReplaceAll(template, "<request_type>", handler.RequestType)
	template = strings.ReplaceAll(template, "<request_type_lower_1st>", util.LowerFirstChar(handler.RequestType))
	template = strings.ReplaceAll(template, "<response_type>", handler.ResponseType)
	template = strings.ReplaceAll(template, "<path>", handler.Path)
	if routing.ApiPublic[handler.Path] {
		template = strings.ReplaceAll(template, "<summary>", "api public")
	} else {
		template = strings.ReplaceAll(template, "<summary>", "permission: "+handler.Permission)
	}

	template = strings.ReplaceAll(template, "<tags>", handler.ServiceName)
	// handle params
	paramDecTemplate, err := util.ReadFileAsString("autogen/template/param_declaration.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingStringTemplate, err := util.ReadFileAsString("autogen/template/param_parsing_string.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingIdTemplate, err := util.ReadFileAsString("autogen/template/param_parsing_id.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramDeclarations := []string{}
	paramParsings := []string{}
	hasIntField := false
	for _, field := range requestType.Fields {

		declaration := paramDecTemplate
		declaration = strings.ReplaceAll(declaration, "<param_name>", field.Name)
		if field.Location == "path" {
			declaration = strings.ReplaceAll(declaration, "<param_location>", field.Location)
			declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(true))
		} else {
			declaration = strings.ReplaceAll(declaration, "<param_location>", "body")
			declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(field.Required))
		}

		field_type := field.Type
		if strings.Contains(field_type, "Request") {
			field_type = handler.ProtoPackage + "." + field_type
		} else {
			switch field_type {
			case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64", "complex64", "complex128", "byte", "rune", "string", "bool", "enum", "bytes":
			case "float", "double":
				field_type = "float64"
			default:
				field_type = handler.ProtoPackage + "." + field_type
			}
		}

		if field.Repeated {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "[]"+field_type)
		} else if field.Location == "enum" {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "string")
		} else if field.Location == "file" {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "file")
		} else {
			declaration = strings.ReplaceAll(declaration, "<param_type>", field_type)
		}
		declaration = strings.ReplaceAll(declaration, "<param_description>", field.Description)
		paramDeclarations = append(paramDeclarations, declaration)

		if field.Location == "path" {
			parsing := paramParsingIdTemplate

			paramContext := "Query"
			if field.Location == "path" {
				paramContext = "Param"
			}

			if strings.HasPrefix(field.Type, "int") {
				hasIntField = true
				parsing = strings.ReplaceAll(parsing, "<param_bit>", strconv.FormatInt(field.Bit, 10))

			} else {
				parsing = paramParsingStringTemplate
			}

			parsing = strings.ReplaceAll(parsing, "<param_name>", field.Name)
			parsing = strings.ReplaceAll(parsing, "<param_context>", paramContext)
			parsing = strings.ReplaceAll(parsing, "<param_name_upper_case_word>", util.ConvertSliceStrToUCWord(strings.Split(field.Name, "_")))
			paramParsings = append(paramParsings, parsing)
		}
	}
	if len(paramDeclarations) > 0 {
		template = strings.ReplaceAll(template, "<param_declarations>", strings.Join(paramDeclarations, "\n"))
	} else {
		template = strings.ReplaceAll(template, "\n<param_declarations>", strings.Join(paramDeclarations, "\n"))
	}
	template = strings.ReplaceAll(template, "<param_parsings>", strings.Join(paramParsings, "\n"))
	if !hasIntField {
		if runtime.GOOS == "windows" {
			regex := regexp.MustCompile(`"strconv"`)
			template = regex.ReplaceAllString(template, "")
		} else {
			regex := regexp.MustCompile(" +\"strconv\"\\n")
			template = regex.ReplaceAllString(template, "")
		}

	}

	err = os.WriteFile(handler.FolderPath+"/"+handlerNameLower1st+".go", []byte(template), 0644)
	if err != nil {
		fmt.Println("failed to write handler file:", err)
		return err
	}

	return nil
}

func generateGetHandler(handler *model.Handler, requestType *model.Message) error {
	template, err := util.ReadFileAsString("autogen/template/handler_get.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	template = strings.ReplaceAll(template, "<proto_folder>", handler.ProtoFolder)
	template = strings.ReplaceAll(template, "<proto_package>", handler.ProtoPackage)
	template = strings.ReplaceAll(template, "<handler_name>", handler.Name)
	handlerNameLower1st := util.LowerFirstChar(handler.Name)
	template = strings.ReplaceAll(template, "<handler_name_lower_1st>", handlerNameLower1st)
	template = strings.ReplaceAll(template, "<request_type>", handler.RequestType)
	template = strings.ReplaceAll(template, "<response_type>", handler.ResponseType)
	template = strings.ReplaceAll(template, "<path>", handler.Path)
	if routing.ApiPublic[handler.Path] {
		template = strings.ReplaceAll(template, "<summary>", "api public")
	} else {
		template = strings.ReplaceAll(template, "<summary>", "permission: "+handler.Permission)
	}
	template = strings.ReplaceAll(template, "<tags>", handler.ServiceName)

	// handle params
	paramDecTemplate, err := util.ReadFileAsString("autogen/template/param_declaration.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingIntTemplate, err := util.ReadFileAsString("autogen/template/param_parsing_int.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingArrIntTemplate, err := util.ReadFileAsString("autogen/template/param_parsing_arr_int.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingStringTemplate, err := util.ReadFileAsString("autogen/template/param_parsing_string.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingArrStringTemplate, err := util.ReadFileAsString("autogen/template/param_parsing_arr_string.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingBoolTemplate, err := util.ReadFileAsString("autogen/template/param_parsing_bool.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingGGTimestampTemplate, err := util.ReadFileAsString("autogen/template/param_parsing_gg_timestamp.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramDeclarations := []string{}
	paramParsings := []string{}
	hasIntField := false
	hasRepeatField := false
	for _, field := range requestType.Fields {
		declaration := paramDecTemplate
		declaration = strings.ReplaceAll(declaration, "<param_name>", field.Name)
		declaration = strings.ReplaceAll(declaration, "<param_location>", field.Location)
		if field.Repeated {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "string")
			declaration = strings.ReplaceAll(declaration, "<param_description>", "value1,value2,value3,...")
		} else {
			declaration = strings.ReplaceAll(declaration, "<param_type>", field.Type)
		}
		if field.Location == "path" {
			declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(true))
		} else {
			declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(field.Required))
		}
		declaration = strings.ReplaceAll(declaration, "<param_description>", field.Description)
		paramDeclarations = append(paramDeclarations, declaration)

		parsing := paramParsingStringTemplate
		if strings.HasPrefix(field.Type, "int") {
			hasIntField = true
			parsing = paramParsingIntTemplate
			if field.Repeated {
				hasRepeatField = true
				parsing = paramParsingArrIntTemplate
			}
		} else if field.Type == "bool" {
			parsing = paramParsingBoolTemplate
		} else if field.Type == "google.protobuf.Timestamp" {
			parsing = paramParsingGGTimestampTemplate
		} else if field.Type == "string" {
			if field.Repeated {
				hasRepeatField = true
				parsing = paramParsingArrStringTemplate
			}
		} else {
			fmt.Println("unsupported param type:", field.Type, ", fallback to string type")
		}

		parsing = strings.ReplaceAll(parsing, "<param_name>", field.Name)
		paramContext := "Query"
		if field.Location == "path" {
			paramContext = "Param"
		}
		parsing = strings.ReplaceAll(parsing, "<param_context>", paramContext)
		parsing = strings.ReplaceAll(parsing, "<param_bit>", strconv.FormatInt(field.Bit, 10))
		parsing = strings.ReplaceAll(parsing, "<param_name_upper_case_word>", util.ConvertSliceStrToUCWord(strings.Split(field.Name, "_")))
		paramParsings = append(paramParsings, parsing)
	}
	if len(paramDeclarations) > 0 {
		template = strings.ReplaceAll(template, "<param_declarations>", strings.Join(paramDeclarations, "\n"))
	} else {
		template = strings.ReplaceAll(template, "\n<param_declarations>", strings.Join(paramDeclarations, "\n"))
	}

	template = strings.ReplaceAll(template, "<param_parsings>", strings.Join(paramParsings, "\n"))
	if !hasIntField {
		if runtime.GOOS == "windows" {
			regex := regexp.MustCompile(`"strconv"`)
			template = regex.ReplaceAllString(template, "")
		} else {
			regex := regexp.MustCompile(" +\"strconv\"\\n")
			template = regex.ReplaceAllString(template, "")
		}
	}
	if !hasRepeatField {
		if runtime.GOOS == "windows" {
			regex := regexp.MustCompile(`"strings"`)
			template = regex.ReplaceAllString(template, "")
		} else {
			regex := regexp.MustCompile(" +\"strings\"\\n")
			template = regex.ReplaceAllString(template, "")
		}
	}

	err = os.WriteFile(handler.FolderPath+"/"+handlerNameLower1st+".go", []byte(template), 0644)
	if err != nil {
		fmt.Println("failed to write handler file:", err)
		return err
	}

	return nil
}

func GenerateRoutingRegistry(handlers []*model.Handler, handlerFolderPaths []string) error {
	fmt.Println("Generating routing registry")
	template, err := util.ReadFileAsString("autogen/template/routing_registry.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}

	handlerImportTemplate, err := util.ReadFileAsString("autogen/template/handler_import.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	handlerImportPaths := []string{}
	for _, path := range handlerFolderPaths {
		handlerImportPath := handlerImportTemplate
		handlerImportPath = strings.ReplaceAll(handlerImportPath, "<handler_folder>", path)
		handlerImportPaths = append(handlerImportPaths, handlerImportPath)
	}
	template = strings.ReplaceAll(template, "<handler_paths>", strings.Join(handlerImportPaths, "\n"))

	routingConfigTemplate, err := util.ReadFileAsString("autogen/template/routing_config.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	routingConfigs := []string{}
	for _, handler := range handlers {
		routingConfig := routingConfigTemplate
		routingConfig = strings.ReplaceAll(routingConfig, "<method>", handler.Method)
		routingConfig = strings.ReplaceAll(routingConfig, "<path>", handler.Path)
		routingConfig = strings.ReplaceAll(routingConfig, "<permission>", handler.Permission)
		routingConfig = strings.ReplaceAll(routingConfig, "<handler_package>", handler.ProtoPackage)
		routingConfig = strings.ReplaceAll(routingConfig, "<handler_name>", handler.Name)
		routingConfig = strings.ReplaceAll(routingConfig, "<service_name>", handler.ServiceName)
		routingConfig = strings.ReplaceAll(routingConfig, "<action>", handler.Name)
		routingConfigs = append(routingConfigs, routingConfig)
	}
	template = strings.ReplaceAll(template, "<routing_configs>", strings.Join(routingConfigs, "\n"))

	err = os.WriteFile(ROUTING_REGISTRY_PATH, []byte(template), 0644)
	if err != nil {
		fmt.Println("failed to write routing registry file:", err)
		return err
	}

	return nil
}
