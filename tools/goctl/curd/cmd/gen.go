package cmd

import (
	"dm.com/toolx/arr"
	"dm.com/toolx/fn/arrfn"
	"errors"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/curd/tpl"
	curdutil "github.com/zeromicro/go-zero/tools/goctl/curd/util"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/gen"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/model"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
	"os"
	"path"
	"strings"

	apiParser "github.com/zeromicro/go-zero/tools/goctl/api/parser"
)

var (
	Gen = cobrax.NewCommand("gen", cobrax.WithRunE(genE))
)

func init() {
	Gen.Flags().StringVar(&api, "api")
	Gen.Flags().StringVar(&dir, "dir")
	Gen.Flags().StringVar(&url, "url")
	Gen.Flags().StringVar(&table, "table")
	Gen.Flags().StringVar(&home, "home")
	Gen.Flags().StringVarWithDefaultValue(&style, "style", config.DefaultFormat)
}

func genE(cmd *cobra.Command, args []string) error {
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}
	if len(api) == 0 {
		return errors.New("missing -api")
	}
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}
	if len(url) == 0 {
		return errors.New("missing -url")
	}
	if len(table) == 0 {
		return errors.New("missing -table")
	}
	return doGenCrud(api, dir, url, table, style)
}

// doGenCrud gen crud files with api file
func doGenCrud(apiFile, dir, url, table, namingStyle string) error {
	apiSpec, err := apiParser.Parse(apiFile)
	if err != nil {
		return err
	}

	tableInfo, err := curdutil.ParseTable(url, table)
	if err != nil {
		return err
	}

	cfg, err := config.NewConfig(namingStyle)
	if err != nil {
		return err
	}
	logx.Must(pathx.MkdirIfNotExist(dir))
	rootPkg, err := golang.GetParentPackage(dir)
	if err != nil {
		return err
	}

	if apiSpec, err = curdutil.ReplaceApi(apiSpec, cfg, tableInfo); err != nil {
		return err
	}
	logx.Must(genModel(dir, cfg, table, tableInfo))
	//logx.Must(genTypes(dir, cfg, apiSpec, tableInfo.Table))
	logx.Must(gogen.GenHandlers(dir, rootPkg, cfg, apiSpec))
	logx.Must(genLogic(dir, rootPkg, cfg, apiSpec, table, tableInfo))
	logx.Must(os.WriteFile(apiFile, []byte(curdutil.ApiSpecToString(apiSpec)), 0666))
	return nil
}

var (
	actionMap = map[string]bool{
		string(tpl.Add):    true,
		string(tpl.Update): true,
		string(tpl.Delete): true,
		string(tpl.Get):    true,
		string(tpl.Page):   true,
	}
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
			action := strings.TrimPrefix(r.Path, "/")
			if r.Curd || actionMap[action] {
				r.Action = action
				err := genLogicByRoute(dir, rootPkg, cfg, g, r, modelName, entityName, typesMap, tableInfo)
				if err != nil {
					return err
				}
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

	colMap := arrfn.ToMap(tableInfo.Columns, func(e *model.Column) (string, *model.Column) {
		return strcase.ToCamel(e.Name), e
	})

	isTime := func(colName string) bool {
		if col, ok := colMap[colName]; ok {
			return col.DataType == "datetime" || col.DataType == "timestamp"
		}
		return false
	}
	isInt64 := func(colName string) bool {
		if col, ok := colMap[colName]; ok {
			return col.DataType == "bigint" || col.DataType == "int" || col.DataType == "tinyint" || strings.Contains(col.DataType, "int")
		}
		return false
	}
	isNull := func(colName string) bool {
		if col, ok := colMap[colName]; ok {
			return col.IsNullAble == "YES"
		}
		return false
	}
	hasDefault := func(colName string) bool {
		if col, ok := colMap[colName]; ok {
			return col.ColumnDefault != nil
		}
		return false
	}

	respItemTypeName := ""
	var respItemMembers []spec.Member
	if route.Action == string(tpl.Page) {
		for i := range respType.Members {
			m := respType.Members[i]
			if m.Name == "List" {
				respItemTypeName = m.Type.Name()[2:]
				respItemMembers = typesMap[respItemTypeName].Members
				break
			}
		}
	}

	return gogen.GenFile(gogen.FileGenConfig{
		Dir:             dir,
		Subdir:          subDir,
		Filename:        goFile + ".go",
		TemplateName:    "logicTemplate",
		Category:        tpl.Category(),
		TemplateFile:    tpl.ActionToLogicFile(route.Action),
		BuiltinTemplate: tpl.Templates[tpl.ActionToLogicFile(route.Action)],
		Data: map[string]any{
			"pkgName":                   subDir[strings.LastIndex(subDir, "/")+1:],
			"imports":                   imports,
			"importModel":               rootPkg + "/internal/model",
			"logic":                     strings.Title(logic),
			"function":                  strings.Title(strings.TrimSuffix(logic, "Logic")),
			"responseType":              responseString,
			"returnString":              returnString,
			"request":                   requestString,
			"route":                     arr.NewMap(group.Annotation.Properties).Get("prefix", "") + route.Path,
			"title":                     arr.NewMap(route.AtDoc.Properties).Get("summary", route.AtDoc.Text),
			"method":                    route.Method,
			"modelName":                 modelName,
			"entityName":                entityName,
			"reqMembers":                reqType.Members,
			"respMembers":               respType.Members,
			"resp":                      "&" + gogen.TypesPacket + "." + respType.Name(),
			"respItemTypeName":          respItemTypeName, // 分页时，返回的数据类型的名称
			"respItemMembers":           respItemMembers,  // 分页时，返回的数据类型的字段
			"lowerStartCamelPrimaryKey": util.EscapeGolangKeyword(stringx.From(stringx.From(tableInfo.PrimaryKey.Name).ToCamel()).Untitle()),
			"primaryKey":                util.EscapeGolangKeyword(stringx.From(tableInfo.PrimaryKey.Name).ToCamel()),
		},
		FuncMap: map[string]any{
			"IsTime": isTime,
			"IsNull": isNull,
			"IsNullTime": func(colName string) bool {
				if hasDefault(colName) {
					return false
				}
				return isTime(colName) && isNull(colName)
			},
			"IsNullInt64": func(colName string) bool {
				if hasDefault(colName) {
					return false
				}
				return isNull(colName) && isInt64(colName)
			},
		},
	})
}

func genTypes(dir string, cfg *config.Config, api *spec.ApiSpec, prefix string) error {
	val, err := gogen.BuildTypes(api.Types)
	if err != nil {
		return err
	}

	typeFilename, err := format.FileNamingFormat(cfg.NamingFormat, "curd_types_"+prefix)
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
