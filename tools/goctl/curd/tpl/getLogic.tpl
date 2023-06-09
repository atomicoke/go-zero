package {{.pkgName}}

import (
    "{{.importModel}}"
    "dm-admin/common/sbuilder"
    "dm-admin/common/errorx"
    "dm.com/toolx/mp"
    "github.com/Masterminds/squirrel"
    "dm-admin/common/globalkey"

	{{.imports}}
)

type {{.logic}} struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func New{{.logic}}(ctx context.Context, svcCtx *svc.ServiceContext) *{{.logic}} {
	return &{{.logic}}{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *{{.logic}}) model() model.{{.modelName}} {
    return l.svcCtx.{{.modelName}}
}

// get sql builder
func (l *{{.logic}}) sql({{.request}}) squirrel.SelectBuilder {
    var (
    	m = l.model()
    	f = m.Fields()
    )
    return sbuilder.Where(m){{- range .reqMembers }}.
       {{ $type := .Type.Name }}
       Eq(f.{{.Name}}, req.{{.Name}}).
   {{- end }}
       Res()
}

/*
@desc  {{.title}}
@route {{.route}}
*/
func (l *{{.logic}}) {{.function}}({{.request}}) {{.responseType}} {
	entity, err := l.model().FindOneByQuery(l.ctx, l.sql(req))
	if err != nil {
		return nil, errorx.Shadow(l, err, {{.title}})
	}

	resp = {{.resp}}{
	{{- range .respMembers }}
	{{ if IsNullTime .Name }}{{.Name}}: mp.NullTimeToString(entity.{{.Name}},globalkey.SysDateFormat),{{else if IsTime .Name}}{{.Name}}: mp.TimeToString(entity.{{.Name}},globalkey.SysDateFormat),{{else if IsNullInt64 .Name}}{{.Name}}: entity.{{.Name}}.Int64,{{else}}{{.Name}}: entity.{{.Name}},{{end}}
    {{- end }}
	}
	{{.returnString}}
}
