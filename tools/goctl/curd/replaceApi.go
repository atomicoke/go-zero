package curd

import (
	"dm.com/toolx/arr"
	"dm.com/toolx/fn/arrfn"
	"dm.com/toolx/fn/mapfn"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/model"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"
	utilformat "github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
	"strings"
)

const prefixKey = "curdPrefix"

func replaceApi(output string, apiSpec *spec.ApiSpec, cfg *config.Config, table *model.Table) error {
	var (
		prefix string
		desc   string
	)

	if len(apiSpec.Info.Properties[prefixKey]) > 0 {
		prefix = strings.Trim(apiSpec.Info.Properties[prefixKey], "\"")
	} else {
		prefix = table.Table
	}
	desc = strings.Trim(apiSpec.Info.Properties["desc"], "\"")
	prefix, err := utilformat.FileNamingFormat(cfg.NamingFormat, prefix)
	if err != nil {
		return err
	}

	group, shouldAppend := findServiceGroup(apiSpec)
	if err := addApiAndType(prefix, desc, table, group, apiSpec); err != nil {
		return err
	}
	if shouldAppend {
		apiSpec.Service.Groups = append(apiSpec.Service.Groups, *group)
	}
	return nil
}

var mapJsonTag = func(name string, comment string) string {
	return fmt.Sprintf("`label:\"%s\" json:\"%s\"`", comment, strcase.ToLowerCamel(name))
}

func addApiAndType(prefix string, desc string, t *model.Table, group *spec.Group, apiSpec *spec.ApiSpec) error {
	var (
		reqType         spec.DefineStruct
		respType        spec.DefineStruct
		addTypeReqName  = stringx.From(prefix).Title() + "AddReq"
		addTypeRespName = stringx.From(prefix).Title() + "AddResp"
		ok              bool
	)
	types := arrfn.ToMap(apiSpec.Types, func(t spec.Type) (string, spec.DefineStruct) {
		return t.Name(), t.(spec.DefineStruct)
	})
	newMembers, pk, err := mapColToMember(t, mapJsonTag, true)
	if err != nil {
		return err
	}

	if reqType, ok = types[addTypeReqName]; !ok {
		apiSpec.Types = append(apiSpec.Types, spec.DefineStruct{
			RawName: addTypeReqName,
			Members: newMembers,
			Docs:    []string{fmt.Sprintf("add %s request", desc)},
		})
	} else {
		mergeMembers(reqType, newMembers)
	}

	if respType, ok = types[addTypeRespName]; !ok {
		apiSpec.Types = append(apiSpec.Types, spec.DefineStruct{
			RawName: addTypeRespName,
			Members: []spec.Member{
				{
					Name: pk.Name.Title(),
					Type: spec.PrimitiveType{RawName: pk.DataType},
					Tag:  mapJsonTag(pk.Name.Source(), pk.Comment),
				},
			},
			Docs: []string{fmt.Sprintf("add %s response", desc)},
		})
	}

	mergeRouters(group, buildRoute(prefix, "post", "add", reqType, respType, "添加 "+desc))

	return nil
}

func mergeRouters(group *spec.Group, router spec.Route) {
	routers := group.Routes
	prevRouterMap := arrfn.ToMap(routers, func(r spec.Route) (string, spec.Route) {
		return r.Path, r
	})

	if _, ok := prevRouterMap[router.Path]; !ok {
		group.Routes = append(routers, router)
	}
}

func mergeMembers(reqType spec.DefineStruct, newMembers []spec.Member) {
	prevMemberMap := arrfn.ToMap(reqType.Members, func(m spec.Member) (string, spec.Member) {
		return m.Name, m
	})

	for i := range newMembers {
		member := newMembers[i]
		if prevMember, ok := prevMemberMap[member.Name]; ok {
			member.Tag = prevMember.Tag
			member.Comment = prevMember.Comment
			member.Docs = prevMember.Docs
		}
	}
	reqType.Members = newMembers
}

func mapColToMember(t *model.Table, mapTag func(name string, comment string) string, skipPri bool) ([]spec.Member, *parser.Primary, error) {
	table, err := parser.ConvertDataType(t, true)
	if err != nil {
		return nil, nil, err
	}

	var members []spec.Member
	for i := range table.Fields {
		field := table.Fields[i]
		if skipPri && field.Name.Source() == table.PrimaryKey.Name.Source() {
			continue
		}
		members = append(members, spec.Member{
			Name: field.Name.ToCamel(),
			Type: spec.PrimitiveType{
				RawName: field.DataType,
			},
			Tag: mapTag(field.Name.Source(), field.Comment),
		})
	}
	return members, &table.PrimaryKey, nil
}

func buildRoute(prefix string, method string, action string, reqType spec.DefineStruct, respType spec.DefineStruct, desc string) spec.Route {
	return spec.Route{
		Method:       method,
		Path:         "/" + action,
		Handler:      stringx.From(action).Title() + stringx.From(prefix).Title() + "Handler",
		AtDoc:        spec.AtDoc{Text: "\"" + desc + "\""},
		RequestType:  reqType,
		ResponseType: respType,
	}
}

func findServiceGroup(apiSpec *spec.ApiSpec) (*spec.Group, bool) {
	var (
		group        *spec.Group
		shouldAppend = false
	)

	arr.Slice(apiSpec.Service.Groups).Find(func(g spec.Group) bool {
		return g.GetAnnotation("curd") == "true"
	}, func(v spec.Group, i int) {
		group = &v

	})
	if group == nil {
		group = &spec.Group{
			Annotation: spec.Annotation{Properties: mapfn.Combine(map[string]string{
				"curd": "true",
			}, apiSpec.Service.Groups[0].Annotation.Properties)},
			Routes: []spec.Route{},
		}
		shouldAppend = true
	}
	return group, shouldAppend
}
