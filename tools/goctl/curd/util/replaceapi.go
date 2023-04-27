package util

import (
	"dm.com/toolx/arr"
	"dm.com/toolx/fn/arrfn"
	"dm.com/toolx/fn/mapfn"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/curd/tpl"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/model"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
	"strings"
)

const prefixKey = "curdPrefix"

func ReplaceApi(apiSpec *spec.ApiSpec, cfg *config.Config, table *model.Table) (*spec.ApiSpec, error) {
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
	prefix = strcase.ToCamel(prefix)

	group := findServiceGroup(apiSpec)
	builder := buildApiAndType(prefix, desc, table, apiSpec)
	addRoute, apiSpec, err := builder(tpl.Add, "添加", "post", members, memberPk)
	if err != nil {
		return nil, err
	}
	updateRoute, apiSpec, err := builder(tpl.Update, "更新", "post", membersAndPk, emptyMembers)
	if err != nil {
		return nil, err
	}
	deleteRoute, apiSpec, err := builder(tpl.Delete, "删除", "post", memberPk, emptyMembers)
	if err != nil {
		return nil, err
	}
	pageRoute, apiSpec, err := builder(tpl.Page, "分页", "get", pageReqMembers, pageRespItemsMembers)
	if err != nil {
		return nil, err
	}
	getRoute, apiSpec, err := builder(tpl.Get, "获取", "get", memberPkGet, membersAndPk)
	if err != nil {
		return nil, err
	}

	mergeRouters(group, addRoute, getRoute, updateRoute, deleteRoute, pageRoute)

	return replaceGroup(group, apiSpec), nil
}

func replaceGroup(group *spec.Group, apiSpec *spec.ApiSpec) *spec.ApiSpec {
	var find = false
	for i := range apiSpec.Service.Groups {
		g := apiSpec.Service.Groups[i]
		if g.GetAnnotation(tpl.Category()) == "true" {
			apiSpec.Service.Groups[i].Routes = group.Routes
			find = true
		}
	}
	if !find {
		apiSpec.Service.Groups = append(apiSpec.Service.Groups, *group)
	}
	return apiSpec
}

func buildApiAndType(prefix string, desc string, t *model.Table, apiSpec *spec.ApiSpec) func(
	action tpl.Action, chinesAction, method string,
	reqMember, respMember memberProvider) (spec.Route, *spec.ApiSpec, error) {
	return func(action tpl.Action, chinesAction, method string, reqMember, respMember memberProvider) (spec.Route, *spec.ApiSpec, error) {
		var (
			actionTitle  = stringx.From(string(action)).Title()
			reqType      spec.DefineStruct
			respType     spec.DefineStruct
			reqTypeName  = strcase.ToCamel(prefix) + actionTitle + "Req"
			respTypeName = strcase.ToCamel(prefix) + actionTitle + "Resp"
			ok           bool
		)
		types := arrfn.ToMap(apiSpec.Types, func(t spec.Type) (string, spec.DefineStruct) {
			return t.Name(), t.(spec.DefineStruct)
		})
		newMembers, pk, err := mapColToMember(t, mapJsonTag, true)
		if err != nil {
			return spec.Route{}, nil, err
		}

		if reqType, ok = types[reqTypeName]; !ok {
			reqType = spec.DefineStruct{
				RawName: reqTypeName,
				Members: reqMember(newMembers, pk),
				Docs:    []string{fmt.Sprintf("%s %s request", action, desc)},
			}
			apiSpec.Types = append(apiSpec.Types, reqType)
		} else {
			reqType.Members = mergeMembers(reqType, reqMember(newMembers, pk))
			for i := range apiSpec.Types {
				if apiSpec.Types[i].Name() == reqTypeName {
					apiSpec.Types[i] = reqType
				}
			}
		}
		if action == tpl.Page {
			respMembers, itemTypeName := pageRespMembers(prefix)
			respType = spec.DefineStruct{
				RawName: respTypeName,
				Members: respMembers,
				Docs:    []string{fmt.Sprintf("%s %s response", action, desc)},
			}
			if _, ok = types[respTypeName]; !ok {
				apiSpec.Types = append(apiSpec.Types, respType)
			} else {
				for i := range apiSpec.Types {
					if apiSpec.Types[i].Name() == respTypeName {
						apiSpec.Types[i] = respType
					}
				}
			}
			itemsMembers := pageRespItemsMembers(newMembers, pk)
			itemsType := spec.DefineStruct{
				RawName: itemTypeName,
			}
			if _, ok = types[itemTypeName]; !ok {
				itemsType = spec.DefineStruct{
					RawName: itemTypeName,
					Members: itemsMembers,
					Docs:    []string{fmt.Sprintf("%s %s response item", action, desc)},
				}
				apiSpec.Types = append(apiSpec.Types, itemsType)
			} else {
				itemsType.Members = mergeMembers(itemsType, itemsMembers)
				for i := range apiSpec.Types {
					if apiSpec.Types[i].Name() == itemTypeName {
						apiSpec.Types[i] = itemsType
					}
				}
			}
		} else {
			if respType, ok = types[respTypeName]; !ok {
				respType = spec.DefineStruct{
					RawName: respTypeName,
					Members: respMember(newMembers, pk),
					Docs:    []string{fmt.Sprintf("%s %s response", action, desc)},
				}
				apiSpec.Types = append(apiSpec.Types, respType)
			} else {
				respType.Members = mergeMembers(respType, respMember(newMembers, pk))
				for i := range apiSpec.Types {
					if apiSpec.Types[i].Name() == respTypeName {
						apiSpec.Types[i] = respType
					}
				}
			}
		}

		return buildRoute(prefix, method, string(action), reqType, respType, chinesAction+" "+desc), apiSpec, nil
	}
}

func mergeRouters(group *spec.Group, newRouter ...spec.Route) {
	routers := group.Routes
	prefix := strings.Trim(group.GetAnnotation("prefix"), "\"")
	prevRouterMap := arrfn.ToMap(routers, func(r spec.Route) (string, spec.Route) {
		return prefix + r.Path, r
	})

	for i := range newRouter {
		router := newRouter[i]
		if _, ok := prevRouterMap[prefix+router.Path]; !ok {
			group.Routes = append(group.Routes, router)
		}
	}
}

func arrRemove[T any](slice []T, i int) []T {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func mergeMembers(reqType spec.DefineStruct, dbMembers []spec.Member) []spec.Member {
	prevMemberMap := arrfn.ToMap(reqType.Members, func(m spec.Member) (string, spec.Member) {
		return m.Name, m
	})

	var copyDbMembers = make([]spec.Member, len(dbMembers))
	var reqTypeMembers = make([]spec.Member, len(reqType.Members))

	copy(copyDbMembers, dbMembers)
	copy(reqTypeMembers, reqType.Members)

	for _, v := range copyDbMembers {
		if _, ok := prevMemberMap[v.Name]; !ok {
			reqTypeMembers = append(reqTypeMembers, v)
		}
	}

	var sortMembers []spec.Member

	for i := 0; i < len(reqTypeMembers); i++ {
		v := reqTypeMembers[i]
		for _, x := range copyDbMembers {
			if v.Name == x.Name {
				sortMembers = append(sortMembers, v)
				reqTypeMembers = arrRemove(reqTypeMembers, i)
				i--
				break
			}
		}
	}

	for _, v := range reqTypeMembers {
		sortMembers = append(sortMembers, v)
	}

	return sortMembers
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
				RawName: maptype(field.DataType),
			},
			Tag: mapTag(field.Name.Source(), field.Comment),
		})
	}
	return members, &table.PrimaryKey, nil
}

func buildRoute(prefix string, method string, action string, reqType spec.DefineStruct, respType spec.DefineStruct, desc string) spec.Route {
	prefix = strcase.ToCamel(prefix)
	route := spec.Route{
		Method:       method,
		Path:         "/" + action,
		Handler:      stringx.From(action).Title() + prefix + "Handler",
		AtDoc:        spec.AtDoc{Text: "\"" + desc + "\""},
		RequestType:  reqType,
		ResponseType: respType,
		Curd:         true,
		Action:       action,
	}
	return route
}

func findServiceGroup(apiSpec *spec.ApiSpec) *spec.Group {
	var (
		group *spec.Group
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
	}
	return group
}
