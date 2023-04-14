func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, session sqlx.Session, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne(ctx, newData.{{.upperStartCamelPrimaryKey}})
	if err!=nil{
		return nil, err
	}
	{{end}}	{{.keys}}
  	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
	query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
	if session != nil{
		return session.ExecCtx(ctx,query, {{.expressionValues}})
	}
	return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
	if session != nil{
		return session.ExecCtx(ctx,query, {{.expressionValues}})
	}
	return m.conn.ExecCtx(ctx, query, {{.expressionValues}}){{end}}
}

// todo 加缓存
func (m *default{{.upperStartCamelObject}}Model) UpdatePart(ctx context.Context, session sqlx.Session, newData *{{.upperStartCamelObject}}Update) (error) {
	keys := arr.Map(&newData.list, func(v KV) string {
		return "`"+v.K + "`=?"
	}).Join(",")
	values := arr.Map(&newData.list, func(v KV) any {
		return v.V
	})

	values.Push(newData.Row.{{ .upperStartCamelPrimaryKey }})

	f := func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = ?", m.table, keys)
		v := values.ToSlice()
		if session != nil {
			return session.ExecCtx(ctx, query, v...)
		}
		return conn.ExecCtx(ctx, query, v...)
	}

	{{ if .withCache }}
	res, err := m.ExecCtx(ctx, f, cacheKey)
	{{else}}
	res, err := f(ctx, m.conn)
	{{end}}
	if err != nil {
		return errors.Wrap(err, "更新失败")
	}

	n, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "获取影响行数")
	}

	if n <= 0 {
		return errors.Errorf("影响行数为:%d", n)
	}

	return nil
}

// 事务
func (m *default{{.upperStartCamelObject}}Model) Trans(ctx context.Context,fn func(ctx context.Context,session sqlx.Session) error) error {
	{{if .withCache}}
	return m.TransactCtx(ctx,func(ctx context.Context,session sqlx.Session) error {
		return  fn(ctx,session)
	})
	{{else}}
	return m.conn.TransactCtx(ctx,func(ctx context.Context,session sqlx.Session) error {
		return  fn(ctx,session)
	})
	{{end}}
}

// 有没有缓存都走这一套
func (m *default{{.upperStartCamelObject}}Model) UpdateVersion(ctx context.Context, session sqlx.Session, oldVersion int64, data *{{.upperStartCamelObject}}) error {
	{{- if .withCache -}}
	cacheKey:= fmt.Sprintf("%s%v", {{.cacheKeyVariable}}, data.{{.upperStartCamelPrimaryKey}})
	{{- end -}}

	values := m.fileExpressionValues(data)
	values = append(values, data.{{.upperStartCamelPrimaryKey}})
	values = append(values, oldVersion)
	f := func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = ? AND version=?", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, values...)
		}
		return conn.ExecCtx(ctx, query, values...)
	}

	{{ if .withCache }}
	res, err := m.ExecCtx(ctx, f, cacheKey)
	{{else}}
	res, err := f(ctx, m.conn)
	{{end}}
	if err != nil {
		return errors.Wrap(err, "更新失败")
	}

	n, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "获取影响行数")
	}

	if n <= 0 {
		return errors.Errorf("影响行数为:%d", n)
	}

	return nil
}
