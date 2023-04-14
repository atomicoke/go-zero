func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context,session sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{.keys}}
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
	query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
	if session != nil{
		return session.ExecCtx(ctx,query,{{.expressionValues}})
	}
	return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}
	query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
	if session != nil{
		return session.ExecCtx(ctx,query,{{.expressionValues}})
	}
	return m.conn.ExecCtx(ctx, query, {{.expressionValues}}){{end}}
}

func (m *default{{.upperStartCamelObject}}Model) InsertPart(ctx context.Context,session sqlx.Session, newData *{{.upperStartCamelObject}}Update) (sql.Result,error) {
	keys := arr.Map(&newData.list, func(v KV) string {
		return "`" + v.K + "`"
	}).Join(",")
	exp := arr.Map(&newData.list, func(v KV) string {
		return "?"
	}).Join(",")
	values := arr.Map(&newData.list, func(v KV) any {
		return v.V
	}).ToSlice()

	{{if .withCache}}{{.keys}}
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
	query := fmt.Sprintf("insert into %s (%s) values (%s)", m.table, keys, exp)
	if session != nil{
		return session.ExecCtx(ctx,query, values...)
	}

	return conn.ExecCtx(ctx, query, values...)
    }, {{.keyValues}}){{else}}
	query := fmt.Sprintf("insert into %s (%s) values (%s)", m.table, keys, exp)
	if session != nil{
		return session.ExecCtx(ctx,query, values ...)
	}
	return m.conn.ExecCtx(ctx, query, values...){{end}}
}

func (m *default{{.upperStartCamelObject}}Model) fileExpressionValues(data *{{.upperStartCamelObject}}) []interface{}{
    return []interface{}{ {{.expressionValues}} }
}

func (m *default{{.upperStartCamelObject}}Model) filedExpression() string {
    return  "{{.expression}}"
}
