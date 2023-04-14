func new{{.upperStartCamelObject}}Model(conn sqlx.SqlConn, c cache.CacheConf) *default{{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
		CachedConn: sqlc.NewConn(conn, c),
		conn:conn,
		table:      {{.table}},
		isCache:{{if .withCache}}true{{else}}false{{end}},
	}
}
