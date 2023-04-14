FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error)
Where(name string, v any) *sbuilder.Find[*{{.upperStartCamelObject}}]
RowBuilder(filed string) squirrel.SelectBuilder
CountBuilder(field string) squirrel.SelectBuilder
SumBuilder(field string) squirrel.SelectBuilder
FindOneByQuery(ctx context.Context,rowBuilder squirrel.SelectBuilder) (*{{.upperStartCamelObject}},error)
FindSum(ctx context.Context,sumBuilder squirrel.SelectBuilder) (float64,error)
FindCount(ctx context.Context,countBuilder squirrel.SelectBuilder) (int64,error)
FindRowsByQuery(ctx context.Context,rowBuilder squirrel.SelectBuilder,orderBy string) ([]*{{.upperStartCamelObject}},error)
FindPageListByPage(ctx context.Context,rowBuilder squirrel.SelectBuilder,page ,pageSize int64,orderBy string) ([]*{{.upperStartCamelObject}},error)
Pagination(ctx context.Context, bulder sbuilder.PaginationBuilder, page int64, limit int64, orderBy string) (list []*{{.upperStartCamelObject}}, total int64, err error)
TableName() string
FiledRows() string
Fields() {{ .upperStartCamelObject }}FieldsType
All(ctx context.Context, builder squirrel.SelectBuilder) (list []*{{.upperStartCamelObject}}, err error)
AllPartial(ctx context.Context, builder squirrel.SelectBuilder) (list []*{{.upperStartCamelObject}}, err error)
FindOnePartial(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*{{.upperStartCamelObject}},error)
FindRowsPartial(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*{{.upperStartCamelObject}},error)
FindPagePartial(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64) ([]*{{.upperStartCamelObject}},error)