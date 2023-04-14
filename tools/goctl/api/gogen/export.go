package gogen

import (
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
)

const (
	Internal      = "internal/"
	TypesPacket   = "types"
	ConfigDir     = internal + "config"
	ContextDir    = internal + "svc"
	HandlerDir    = internal + "handler"
	LogicDir      = internal + "logic"
	MiddlewareDir = internal + "middleware"
	TypesDir      = internal + typesPacket
	GroupProperty = "group"
)

func GetLogicName(route spec.Route) string {
	return getLogicName(route)
}

func GenHandlers(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	return genHandlers(dir, rootPkg, cfg, api)
}

func GenLogicImports(route spec.Route, parentPkg string) string {
	return genLogicImports(route, parentPkg)
}

func ResponseGoTypeName(r spec.Route, pkg ...string) string {
	return responseGoTypeName(r, pkg...)
}

func RequestGoTypeName(r spec.Route, pkg ...string) string {
	return requestGoTypeName(r, pkg...)
}
func GetLogicFolderPath(group spec.Group, route spec.Route) string {
	return getLogicFolderPath(group, route)
}

func GenFile(c FileGenConfig) error {
	return genFile(fileGenConfig{
		dir:             c.Dir,
		subdir:          c.Subdir,
		filename:        c.Filename,
		templateName:    c.TemplateName,
		category:        c.Category,
		templateFile:    c.TemplateFile,
		builtinTemplate: c.BuiltinTemplate,
		data:            c.Data,
	})
}

type FileGenConfig struct {
	Dir             string
	Subdir          string
	Filename        string
	TemplateName    string
	Category        string
	TemplateFile    string
	BuiltinTemplate string
	Data            any
}
