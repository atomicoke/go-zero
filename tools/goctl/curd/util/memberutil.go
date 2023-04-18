package util

import (
	"github.com/iancoleman/strcase"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"
)

type memberProvider func(members []spec.Member, pk *parser.Primary) []spec.Member

func memberPk(members []spec.Member, pk *parser.Primary) []spec.Member {
	return []spec.Member{
		{
			Name: pk.Name.Title(),
			Type: spec.PrimitiveType{RawName: pk.DataType},
			Tag:  mapJsonTag(pk.Name.Source(), pk.Comment),
		},
	}
}

func memberPkGet(members []spec.Member, pk *parser.Primary) []spec.Member {
	return []spec.Member{
		{
			Name: pk.Name.Title(),
			Type: spec.PrimitiveType{RawName: pk.DataType},
			Tag:  mapFormTag(pk.Name.Source(), pk.Comment),
		},
	}
}

func members(members []spec.Member, pk *parser.Primary) []spec.Member {
	return members
}

func membersAndPk(members []spec.Member, pk *parser.Primary) []spec.Member {
	return append(memberPk(members, pk), members...)
}

func emptyMembers(members []spec.Member, pk *parser.Primary) []spec.Member {
	return []spec.Member{}
}

func pageReqMembers(members []spec.Member, pk *parser.Primary) []spec.Member {
	return append([]spec.Member{
		{
			Name: "Page",
			Type: spec.PrimitiveType{RawName: "int64"},
			Tag:  mapFormTagWithValid("page", "页码", "number,gte=1"),
		},
		{
			Name: "Limit",
			Type: spec.PrimitiveType{RawName: "int64"},
			Tag:  mapFormTagWithValid("limit", "每页数量", "number,gte=1,lte=100"),
		},
	}, members...)
}

func pageRespItemsMembers(members []spec.Member, pk *parser.Primary) []spec.Member {
	return membersAndPk(members, pk)
}

/*

 {{ if eq .Type.Name "string"  }}
            Eq(m.Fields().{{.Name}},req.{{.Name}}).
       {{end}}
       {{ else if eq .Type.Name "int64"  }}
            EqOn(m.Fields().{{.Name}} != 0,m.Fields().{{.Name}},req.{{.Name}}).
       {{end}}

*/

func pageRespMembers(prefix string) ([]spec.Member, string) {
	prefix = strcase.ToCamel(prefix)
	name := prefix + "PageItem"
	return append([]spec.Member{
		{
			Name: "Page",
			Type: spec.PrimitiveType{RawName: "int64"},
			Tag:  mapJsonTag("page", ""),
		},
		{
			Name: "Limit",
			Type: spec.PrimitiveType{RawName: "int64"},
			Tag:  mapJsonTag("limit", "每页数量"),
		},
		{
			Name: "Total",
			Type: spec.PrimitiveType{RawName: "int64"},
			Tag:  mapJsonTag("total", "总数"),
		},
	}, spec.Member{Name: "List", Type: spec.DefineStruct{RawName: "[]" + name}, Tag: mapJsonTag("list", "列表数据")}), name
}
