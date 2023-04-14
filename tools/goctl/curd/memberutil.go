package curd

import (
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

func members(members []spec.Member, pk *parser.Primary) []spec.Member {
	return members
}

func membersAndPk(members []spec.Member, pk *parser.Primary) []spec.Member {
	return append(memberPk(members, pk), members...)
}

func emptyMembers(members []spec.Member, pk *parser.Primary) []spec.Member {
	return []spec.Member{}
}
