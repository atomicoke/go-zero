package {{.pkgName}}

import (
    "{{.importModel}}"
    "dm-admin/common/errorx"
    "dm-admin/common/sbuilder"

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

// update sql builder
func (l *{{.logic}}) sql({{.request}}) *sbuilder.UpdateSql {
    var (
    	m = l.model()
    	f = m.Fields()
    )
    return sbuilder.Update(m){{- range .reqMembers }}.
       {{ $type := .Type.Name }}
       Eq(f.{{.Name}}, req.{{.Name}})
   {{- end }}
}

/*
@desc  {{.title}}
@route {{.route}}
*/
func (l *{{.logic}}) {{.function}}({{.request}}) {{.responseType}} {
	err = l.model().UpdateCtxWithBuilder(l.ctx, l.sql(req))
	if err != nil {
		return nil, errorx.Shadow(l, err, {{.title}})
	}

	resp = {{.resp}}{}
	{{.returnString}}
}
