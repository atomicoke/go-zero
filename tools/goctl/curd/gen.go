package curd

import (
	"dm.com/toolx/arr"
	"dm.com/toolx/fn/arrfn"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/gen"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/model"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"os"
	"path"
	"strings"
)

func genLogic(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec, tableName string, tableInfo *model.Table) error {
	modelName := strcase.ToCamel(tableName + "Model")
	entityName := strcase.ToCamel(tableName)
	typesMap := arrfn.ToMap(api.Types, func(e spec.Type) (string, spec.DefineStruct) {
		return e.Name(), e.(spec.DefineStruct)
	})
	for _, g := range api.Service.Groups {
		if g.GetAnnotation("curd") != "true" {
			continue
		}

		for _, r := range g.Routes {
			if !r.Curd {
				continue
			}
			err := genLogicByRoute(dir, rootPkg, cfg, g, r, modelName, entityName, typesMap, tableInfo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func genModel(dir string, cfg *config.Config, tableName string, table *model.Table) error {
	dir = path.Join(dir, "internal", "model")
	generator, err := gen.NewDefaultGenerator(dir, cfg,
		gen.WithConsoleOption(console.NewConsole(true)))
	if err != nil {
		return err
	}

	return generator.StartFromInformationSchema(map[string]*model.Table{tableName: table}, false, false)
}

func genLogicByRoute(dir, rootPkg string, cfg *config.Config, group spec.Group, route spec.Route,
	modelName, entityName string, typesMap map[string]spec.DefineStruct, tableInfo *model.Table) error {
	logic := gogen.GetLogicName(route)
	goFile, err := format.FileNamingFormat(cfg.NamingFormat, logic)
	if err != nil {
		return err
	}

	imports := gogen.GenLogicImports(route, rootPkg)
	var (
		responseString string
		returnString   string
		requestString  string
		reqType        = typesMap[route.RequestTypeName()]
		respType       = typesMap[route.ResponseTypeName()]
	)

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
		Data: map[string]any{
			"pkgName":      subDir[strings.LastIndex(subDir, "/")+1:],
			"imports":      imports,
			"importModel":  rootPkg + "/internal/model",
			"logic":        strings.Title(logic),
			"function":     strings.Title(strings.TrimSuffix(logic, "Logic")),
			"responseType": responseString,
			"returnString": returnString,
			"request":      requestString,
			"route":        arr.NewMap(group.Annotation.Properties).Get("prefix", "") + route.Path,
			"title":        arr.NewMap(route.AtDoc.Properties).Get("summary", route.AtDoc.Text),
			"method":       route.Method,
			"modelName":    modelName,
			"entityName":   entityName,
			"reqMembers":   reqType.Members,
			"respMembers":  respType.Members,
			"resp":         "&" + gogen.TypesPacket + "." + respType.Name(),
		},
		FuncMap: map[string]any{
			"IsTime": func(colName string) bool {
				for i := range tableInfo.Columns {
					if strcase.ToCamel(tableInfo.Columns[i].Name) == colName {
						b := tableInfo.Columns[i].DataType == "datetime"
						if b {
							fmt.Println("is time: ", colName)
						}
						return b
					}
				}
				return false
			},
		},
	})
}

func genTypes(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	val, err := gogen.BuildTypes(api.Types)
	if err != nil {
		return err
	}

	typeFilename, err := format.FileNamingFormat(cfg.NamingFormat, "curd_types")
	if err != nil {
		return err
	}

	typeFilename = typeFilename + ".go"
	filename := path.Join(dir, gogen.TypesDir, typeFilename)
	_ = os.Remove(filename)

	return gogen.GenFile(gogen.FileGenConfig{
		Dir:             dir,
		Subdir:          gogen.TypesDir,
		Filename:        typeFilename,
		TemplateName:    "typesTemplate",
		Category:        gogen.CategoryE,
		TemplateFile:    gogen.TypesTemplateFile,
		BuiltinTemplate: gogen.TypesTemplate,
		Data: map[string]any{
			"types":        val,
			"containsTime": false,
		},
	})
}
