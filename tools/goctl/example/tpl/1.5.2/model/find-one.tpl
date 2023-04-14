func (m *default{{.upperStartCamelObject}}Model) FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.QueryRowCtx(ctx, &resp, {{.cacheKeyVariable}}, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query :=  fmt.Sprintf("select %s from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}} limit 1", {{.lowerStartCamelObject}}Rows, m.table)
		return conn.QueryRowCtx(ctx, v, query, {{.lowerStartCamelPrimaryKey}})
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{else}}query := fmt.Sprintf("select %s from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}} limit 1", {{.lowerStartCamelObject}}Rows, m.table)
	var resp {{.upperStartCamelObject}}
	err := m.conn.QueryRowCtx(ctx, &resp, query, {{.lowerStartCamelPrimaryKey}})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{end}}
}

func (m *default{{.upperStartCamelObject}}Model) Where(name string, v any) *sbuilder.Find[*{{.upperStartCamelObject}}] {
	return sbuilder.NewFind[*{{.upperStartCamelObject}}](m).Eq(name, v)
}

func (m *default{{.upperStartCamelObject}}Model) FindOneByQuery(ctx context.Context,rowBuilder squirrel.SelectBuilder) (*{{.upperStartCamelObject}},error) {
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp {{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowNoCacheCtx(ctx,&resp, query, values...){{else}}
	err = m.conn.QueryRowCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}


func (m *default{{.upperStartCamelObject}}Model) FindSum(ctx context.Context,sumBuilder squirrel.SelectBuilder) (float64,error) {
	query, values, err := sumBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var resp float64
	{{if .withCache}}err = m.QueryRowNoCacheCtx(ctx,&resp, query, values...){{else}}
	err = m.conn.QueryRowCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindCount(ctx context.Context,countBuilder squirrel.SelectBuilder) (int64,error) {
	query, values, err := countBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var resp int64
	{{if .withCache}}err = m.QueryRowNoCacheCtx(ctx,&resp, query, values...){{else}}
	err = m.conn.QueryRowCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindRowsByQuery(ctx context.Context,rowBuilder squirrel.SelectBuilder,orderBy string) ([]*{{.upperStartCamelObject}},error) {
	if orderBy == ""{
		rowBuilder = rowBuilder.OrderBy("id DESC")
	}else{
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*{{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowsNoCacheCtx(ctx,&resp, query, values...){{else}}
	err = m.conn.QueryRowsCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindPageListByPage(ctx context.Context,rowBuilder squirrel.SelectBuilder,page ,pageSize int64,orderBy string) ([]*{{.upperStartCamelObject}},error) {
	if orderBy == ""{
		rowBuilder = rowBuilder.OrderBy("id DESC")
	}else{
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	if page < 1{
		page = 1
	}

	offset := (page - 1) * pageSize

	query, values, err := rowBuilder.Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*{{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowsNoCacheCtx(ctx,&resp, query, values...){{else}}
	err = m.conn.QueryRowsCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// export logic
func (m *default{{.upperStartCamelObject}}Model) RowBuilder(filed string) squirrel.SelectBuilder {
	if filed == "*" || filed == ""{
		return squirrel.Select({{.lowerStartCamelObject}}Rows).From(m.table)
	}
	return squirrel.Select(filed).From(m.table)
}

// export logic
func (m *default{{.upperStartCamelObject}}Model) CountBuilder(field string) squirrel.SelectBuilder {
	if field == "" {
        field = "*"
	}
	return squirrel.Select("COUNT("+field+")").From(m.table)
}

// export logic
func (m *default{{.upperStartCamelObject}}Model) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM("+field+"),0)").From(m.table)
}

// 简单分页方法
func (m *default{{.upperStartCamelObject}}Model) Pagination(ctx context.Context, builder sbuilder.PaginationBuilder, page int64, limit int64, orderBy string) (list []*{{.upperStartCamelObject}}, total int64, err error) {
	var where, count = builder.Res()
	if orderBy == "" {
		where = where.OrderBy("id DESC")
	} else {
		where = where.OrderBy(orderBy)
	}

	if page < 1 {
		page = 1
	}
	if limit == 0{
        limit = 10
	}

	offset := (page - 1) * limit

	{
		query, values, err := where.Offset(uint64(offset)).Limit(uint64(limit)).ToSql()
		if err != nil {
			return nil, 0, err
		}

        {{if .withCache}}err = m.QueryRowsNoCacheCtx(ctx,&list, query, values...){{else}}
        err = m.conn.QueryRowsCtx(ctx,&list, query, values...)
        {{end}}
		if err != nil {
			return nil, 0, err
		}
	}
	{
		query, values, err := count.ToSql()
		if err != nil {
			return nil, 0, err
		}

        {{if .withCache}}err = m.QueryRowNoCacheCtx(ctx,&total, query, values...){{else}}
        err = m.conn.QueryRowCtx(ctx,&total, query, values...)
        {{end}}
        
		if err != nil {
			return nil, 0, err
		}
	}

	return list, total, nil
}

func (m *default{{.upperStartCamelObject}}Model) TableName() string {
	return m.tableName()
}

func (m *default{{.upperStartCamelObject}}Model) FiledRows() string {
	return {{.lowerStartCamelObject}}Rows
}

func (m *default{{.upperStartCamelObject}}Model) Fields() {{ .upperStartCamelObject }}FieldsType {
	return {{ .upperStartCamelObject }}Fields
}

func (m *default{{.upperStartCamelObject}}Model) All(ctx context.Context, builder squirrel.SelectBuilder) (list []*{{.upperStartCamelObject}}, err error) {
	query, values, err := builder.ToSql()
	if err != nil {
		return nil, errorx.Wrapf(err, "builder_to_sql:%v", m.tableName())
	}

	{{if .withCache}}err = m.QueryRowsNoCacheCtx(ctx,&list, query, values...){{else}}
	err = m.conn.QueryRowsCtx(ctx,&list, query, values...)
	{{end}}
	if err != nil {
		return nil, errorx.Wrapf(err, "query_rows:%v", m.tableName())
	}

	return list, nil
}

func (m *default{{.upperStartCamelObject}}Model) AllPartial(ctx context.Context, builder squirrel.SelectBuilder) (list []*{{.upperStartCamelObject}}, err error) {
	query, values, err := builder.ToSql()
	if err != nil {
		return nil, errorx.Wrapf(err, "builder_to_sql:%v", m.tableName())
	}

	err = m.conn.QueryRowsPartialCtx(ctx, &list, query, values...)

	if err != nil {
		return nil, errorx.Wrapf(err, "query_rows:%v", m.tableName())
	}

	return list, nil
}

func (m *default{{.upperStartCamelObject}}Model) FindOnePartial(ctx context.Context,rowBuilder squirrel.SelectBuilder) (*{{.upperStartCamelObject}},error) {
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp {{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowsNoCacheCtx(ctx,&resp, query, values...){{else}}
	err = m.conn.QueryRowPartialCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindRowsPartial(ctx context.Context,rowBuilder squirrel.SelectBuilder) ([]*{{.upperStartCamelObject}},error) {
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*{{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowsNoCacheCtx(ctx,&resp, query, values...){{else}}
	err = m.conn.QueryRowsPartialCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindPagePartial(ctx context.Context,rowBuilder squirrel.SelectBuilder,page ,pageSize int64) ([]*{{.upperStartCamelObject}},error) {
	if page < 1{
		page = 1
	}

	offset := (page - 1) * pageSize

	query, values, err := rowBuilder.Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*{{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowsNoCacheCtx(ctx,&resp, query, values...){{else}}
	err = m.conn.QueryRowsPartialCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}


{{if not .withCache}}
func (m *default{{.upperStartCamelObject}}Model) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", "", primary)
}
func (m *default{{.upperStartCamelObject}}Model) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where {{.originalPrimaryField}} = {{if .postgreSql}}$1{{else}}?{{end}} limit 1", {{.lowerStartCamelObject}}Rows, m.table )
	return conn.QueryRowCtx(ctx, v, query, primary)
}
{{end}}
