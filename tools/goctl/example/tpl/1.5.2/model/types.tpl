type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
	}

	default{{.upperStartCamelObject}}Model struct {
		sqlc.CachedConn
		conn sqlx.SqlConn
		table string
        isCache bool
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}

	{{.upperStartCamelObject}}Update struct {
        Row *{{.upperStartCamelObject}}
		list arr.Vector[KV]
	}
)

type  {{ .upperStartCamelObject }}FieldsType  = struct {
{{- range .tableFields}}
	{{ .NameHump }} string
{{- end}}
}

var  {{ .upperStartCamelObject }}Fields  = {{ .upperStartCamelObject }}FieldsType {
{{- range .tableFields}}
	{{ .NameHump }} : "{{ .NameOriginal }}",
{{- end}}
}

{{- range .tableFields}}
func (r *{{ $.upperStartCamelObject }}Update ) Set{{.NameHump}}(v {{.DataType}}) *{{ $.upperStartCamelObject }}Update {
    r.Row.{{.NameHump}} = v
	r.list = append(r.list, KV{"{{ .NameOriginal }}", v})
	return r
}
{{- end}}

func (r *{{ $.upperStartCamelObject }}) ToEntity() *{{ $.upperStartCamelObject }}Update {
    return &{{ $.upperStartCamelObject }}Update{
        Row : r,
    }
}

