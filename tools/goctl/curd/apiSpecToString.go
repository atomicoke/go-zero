package curd

import (
	"dm.com/toolx/arr"
	"fmt"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"sort"
	"strings"
)

var (
	groupKey      = []string{"jwt", "group", "prefix", "middleware", "curd"}
	sortServerKey = map[string]int{
		"jwt":        0,
		"group":      1,
		"prefix":     2,
		"middleware": 3,
		"curd":       4,
	}
	sortInfoKey = map[string]int{
		"title":      0,
		"desc":       1,
		"author":     2,
		"email":      3,
		"curdPrefix": 4,
	}
)

func apiSpecToString(apiSpec *spec.ApiSpec) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("syntax = %s\n\n", apiSpec.Syntax.Version))

	if len(apiSpec.Info.Properties) > 0 {
		writeInfo(sb, apiSpec.Info)
		sb.WriteString("\n")
	}

	if len(apiSpec.Imports) > 0 {
		writeImports(sb, apiSpec.Imports)
		sb.WriteString("\n")
	}

	writeServices(sb, apiSpec.Service)
	writeTypes(sb, apiSpec.Types)
	return sb.String()
}

func writeTypes(sb *strings.Builder, types []spec.Type) {
	if len(types) == 0 {
		return
	}
	sb.WriteString("type (\n")
	for _, v := range types {
		writeType(sb, v)
	}
	sb.WriteString(")\n")
}

func writeType(sb *strings.Builder, t spec.Type) {
	switch t.(type) {
	case spec.DefineStruct:
		writeStruct(sb, t.(spec.DefineStruct))
	default:
		return
	}
}

func writeStruct(sb *strings.Builder, defineStruct spec.DefineStruct) {
	sb.WriteString(fmt.Sprintf("\t%s {\n", defineStruct.Name()))
	for i := range defineStruct.Members {
		member := defineStruct.Members[i]
		writeDoc(sb, member.Docs)
		sb.WriteString(fmt.Sprintf("\t\t%s %s %s %s", member.Name, member.Type.Name(), member.Tag, member.Comment))
		sb.WriteString("\n")
	}
	sb.WriteString("\t}\n\n")
}

func writeInfo(sb *strings.Builder, info spec.Info) {
	if len(info.Properties) == 0 {
		return
	}
	sb.WriteString("info (\n")
	writeAnnotation2(sb, info.Properties, sortInfoKey)
	sb.WriteString(")\n")
}

func writeImports(sb *strings.Builder, imports []spec.Import) {
	if len(imports) == 0 {
		return
	}
	if len(imports) == 1 {
		item := imports[0]
		writeDoc(sb, item.Doc)
		sb.WriteString(fmt.Sprintf("import \"%s\";", item.Value))
		writeComment(sb, item.Comment)
		return
	}

	sb.WriteString("import (\n")
	for _, v := range imports {
		writeDoc(sb, v.Doc)
		sb.WriteString(fmt.Sprintf("\t\"%s\";", v))
		writeComment(sb, v.Comment)
	}
	sb.WriteString(")\n")
}

func writeServices(sb *strings.Builder, service spec.Service) {
	servicesMap := groupByProperties(service, groupKey)
	name := service.Name
	if len(servicesMap) == 0 {
		return
	}
	for k := range servicesMap {
		services := servicesMap[k]
		services.Each(func(group spec.Group) {
			writeServer(sb, group)
			writeService(sb, group, name)
			sb.WriteString("\n")
		})
	}
}

func writeService(sb *strings.Builder, group spec.Group, name string) {
	if len(group.Routes) == 0 {
		return
	}
	sb.WriteString(fmt.Sprintf("service %s {\n", name))
	writeRoutes(sb, group.Routes)
	sb.WriteString("}\n")
}

func writeRoutes(sb *strings.Builder, routes []spec.Route) {
	for _, v := range routes {
		writeRoute(sb, v)
		sb.WriteString("\n")
	}
}

func writeRoute(sb *strings.Builder, route spec.Route) {
	if len(route.AtDoc.Text) > 0 {
		sb.WriteString(fmt.Sprintf("\t@doc %s\n", route.AtDoc.Text))
	} else {
		sb.WriteString(fmt.Sprintf("\t@doc(\n"))
		writeAnnotation3(sb, route.AtDoc.Properties)
		sb.WriteString(fmt.Sprintf("\t)\n"))
	}

	if len(route.HandlerDoc) > 0 {
		sb.WriteString(route.HandlerDoc[0])
		sb.WriteString("\n")
	}
	sb.WriteString(fmt.Sprintf("\t@handler %s", route.Handler))
	writeComment(sb, route.HandlerComment)

	handleUrl := fmt.Sprintf("\t%s %s", route.Method, route.Path)
	if len(route.RequestTypeName()) > 0 {
		handleUrl = fmt.Sprintf("%s (%s)", handleUrl, route.RequestTypeName())
	}
	if len(route.ResponseTypeName()) > 0 {
		handleUrl = fmt.Sprintf("%s returns (%s)", handleUrl, route.ResponseTypeName())
	}
	writeDoc(sb, route.Doc)
	sb.WriteString(handleUrl)
	writeComment(sb, route.Comment)
}

func writeServer(sb *strings.Builder, group spec.Group) {
	if len(group.Annotation.Properties) == 0 {
		return
	}
	sb.WriteString("@server (\n")
	writeAnnotation2(sb, group.Annotation.Properties, sortServerKey)
	sb.WriteString(")\n")
}

func groupByProperties(service spec.Service, key []string) map[string]*arr.Vector[spec.Group] {
	return arr.Slice(service.Groups).Group(func(item spec.Group) string {
		if len(item.Annotation.Properties) == 0 {
			return ""
		}
		var byKey = ""
		for _, v := range key {
			if val, ok := item.Annotation.Properties[v]; ok {
				byKey += fmt.Sprintf("%s:%s", v, val)
				break
			}
		}
		return byKey
	}).ToMap()
}

func writeAnnotation(sb *strings.Builder, p map[string]string) {
	for k, v := range p {
		sb.WriteString(fmt.Sprintf("\t%s: %s\n", k, v))
	}
}

func writeAnnotation2(sb *strings.Builder, p map[string]string, sortServerKey map[string]int) {
	//sort map
	var keys []string
	for k := range p {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return sortServerKey[keys[i]] < sortServerKey[keys[j]]
	})
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("\t%s: %s\n", k, p[k]))
	}
}

func writeAnnotation3(sb *strings.Builder, p map[string]string) {
	for k, v := range p {
		sb.WriteString(fmt.Sprintf("\t\t%s: %s\n", k, v))
	}
}

func writeDoc(sb *strings.Builder, doc spec.Doc) {
	if len(doc) > 0 {
		sb.WriteString(doc[0])
	}
}

func writeComment(sb *strings.Builder, comment spec.Doc) {
	if len(comment) > 0 {
		sb.WriteString(comment[0])
	}
	sb.WriteString("  \n")
}
