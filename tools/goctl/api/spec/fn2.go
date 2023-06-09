package spec

var timeMap = map[string]bool{
	"Time":        true,
	"time":        true,
	"At":          true,
	"at":          true,
	"Date":        true,
	"date":        true,
	"created_at":  true,
	"updated_at":  true,
	"deleted_at":  true,
	"createAt":    true,
	"updateAt":    true,
	"deleteAt":    true,
	"createdAt":   true,
	"updatedAt":   true,
	"deletedAt":   true,
	"CreateTime":  true,
	"UpdateTime":  true,
	"DeleteTime":  true,
	"DeletedTime": true,
}

var timeTypeMap = map[string]bool{
	"time.Time":    true,
	"sql.NullTime": true,
}

func (m Member) IsTime() bool {
	return timeMap[m.Name]
}

func (m Member) IsTimeType() bool {
	return timeTypeMap[m.Type.Name()]
}
