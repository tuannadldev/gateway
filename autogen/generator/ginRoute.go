package generator

import (
	"fmt"
	"gateway/autogen/model"
	"gateway/autogen/util"
	"os"
	"strings"
)

const GIN_ROUTES_PATH = "application/routing/delivery/ginRoutes.go"

func GenerateGinRoutes(protos []*model.Proto) error {
	fmt.Println("Generating gin routes")
	template, err := util.ReadFileAsString("autogen/template/gin_routes.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}

	groupTemplate, err := util.ReadFileAsString("autogen/template/route_group.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	entryTemplate, err := util.ReadFileAsString("autogen/template/route_entry.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	routes := []string{}
	for _, proto := range protos {
		for _, service := range proto.Services {
			if service.Name == "LoginService" {
				continue // skip login service routes, write manually in controller.go
			}
			group := groupTemplate
			group = strings.ReplaceAll(group, "<group>", service.Path)
			routes = append(routes, group)
			for _, api := range service.Apis {
				entry := entryTemplate
				entry = strings.ReplaceAll(entry, "<method>", api.Method)
				entry = strings.ReplaceAll(entry, "<path>", api.Path)
				routes = append(routes, entry)
			}
			routes = append(routes, "") // group separator line
		}
	}
	template = strings.ReplaceAll(template, "<routes>", strings.Join(routes, "\n"))

	err = os.WriteFile(GIN_ROUTES_PATH, []byte(template), 0644)
	if err != nil {
		fmt.Println("failed to write gin routes file:", err)
		return err
	}

	return nil
}
