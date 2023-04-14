package curd

import (
	"dm.com/toolx/arr"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"strings"
)

func genLogic(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	for _, g := range api.Service.Groups {
		if g.GetAnnotation("curd") != "true" {
			continue
		}

		for _, r := range g.Routes {
			if !r.Curd {
				continue
			}
			err := genLogicByRoute(dir, rootPkg, cfg, g, r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func genLogicByRoute(dir, rootPkg string, cfg *config.Config, group spec.Group, route spec.Route) error {
	logic := gogen.GetLogicName(route)
	goFile, err := format.FileNamingFormat(cfg.NamingFormat, logic)
	if err != nil {
		return err
	}

	imports := gogen.GenLogicImports(route, rootPkg)
	var responseString string
	var returnString string
	var requestString string
	if len(route.ResponseTypeName()) > 0 {
		resp := gogen.ResponseGoTypeName(route, gogen.TypesPacket)
		responseString = "(resp " + resp + ", err error)"
		returnString = "return"
	} else {
		responseString = "error"
		returnString = "return nil"
	}
	if len(route.RequestTypeName()) > 0 {
		requestString = "req *" + gogen.RequestGoTypeName(route, gogen.TypesPacket)
	}

	subDir := gogen.GetLogicFolderPath(group, route)
	return gogen.GenFile(gogen.FileGenConfig{
		Dir:             dir,
		Subdir:          subDir,
		Filename:        goFile + ".go",
		TemplateName:    "logicTemplate",
		Category:        category,
		TemplateFile:    actionToLogicFile(route.Action),
		BuiltinTemplate: templates[actionToLogicFile(route.Action)],
		Data: map[string]string{
			"pkgName":      subDir[strings.LastIndex(subDir, "/")+1:],
			"imports":      imports,
			"logic":        strings.Title(logic),
			"function":     strings.Title(strings.TrimSuffix(logic, "Logic")),
			"responseType": responseString,
			"returnString": returnString,
			"request":      requestString,
			"route":        arr.NewMap(group.Annotation.Properties).Get("prefix", "") + route.Path,
			"title":        arr.NewMap(route.AtDoc.Properties).Get("summary", ""),
			"method":       route.Method,
		},
	})
}
