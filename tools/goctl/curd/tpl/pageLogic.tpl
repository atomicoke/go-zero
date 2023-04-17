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

// page sql builder
func (l *{{.logic}}) sql({{.request}}) *sbuilder.Page {
    var (
    	m = l.model()
    	f = m.Fields()
    )
    return sbuilder.BuildPage("",m){{- range .reqMembers }}{{ if eq .Name "Page"}}{{ else if eq .Name "Limit"}}{{else}}.
    Eq(f.{{.Name}}, req.{{.Name}}){{end}}{{- end }}
}

const orderBy = "{{.lowerStartCamelPrimaryKey}} DESC"

/*
@desc  {{.title}}
@route {{.route}}
*/
func (l *{{.logic}}) {{.function}}({{.request}}) {{.responseType}} {
	list, total, err := l.model().Pagination(l.ctx, l.sql(req), req.Page, req.Limit,orderBy)
	if err != nil {
		return nil, errorx.Shadow(l, err, {{.title}})
	}

	resp = {{.resp}}{
		Total: total,
    	Page:  req.Page,
    	Limit: req.Limit,
    	List:  mapList(list),
	}
	{{.returnString}}
}

func mapList(list []*model.{{.modelName}}) []types.{{.respItemTypeName}} {
    var resp []types.{{.respItemTypeName}}
    for _, item := range list {
        v := types.{{.respItemTypeName}}{
            
        }
        resp = append(resp, v)
    }
    return resp
}
