package {{.pkgName}}

import (
    "{{.importModel}}"
    "dm-admin/common/sbuilder"
    "dm-admin/common/errorx"
    "dm-admin/common/utils"
    "github.com/Masterminds/squirrel"

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

// build model.{{.modelName}}
func (l *{{.logic}}) buildEntity({{.request}}) model.{{.entityName}} {
    return model.TestThinkRechargeV2{ {{- range .reqMembers }}
        {{ if IsTime .Name }}{{.Name}}: utils.MapTime(req.{{.Name}}),{{else}}{{.Name}}: req.{{.Name}},{{end}}
   {{- end }}
    }
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
	    {{ if .IsTime }}{{.Name}}: utils.MapTime(entity.{{.Name}}),{{else}}{{.Name}}: entity.{{.Name}},{{end}}
    {{- end }}
	}
	{{.returnString}}
}
