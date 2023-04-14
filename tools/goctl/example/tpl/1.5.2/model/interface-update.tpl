Update(ctx context.Context,session sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result, error)
Trans(ctx context.Context,fn func(context context.Context,session sqlx.Session) error) error
UpdateVersion(ctx context.Context, session sqlx.Session, oldVersion int64, data *{{.upperStartCamelObject}}) error
UpdatePart(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}Update) error
