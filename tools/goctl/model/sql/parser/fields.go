package parser

type TableFieldFull struct {
	*Field
	NameHump string
}

func ToTableFieldFull(field *Field) TableFieldFull {
	return TableFieldFull{
		Field:    field,
		NameHump: field.Name.ToCamel(),
	}
}

func ConvertFieldToFull(field []*Field) []TableFieldFull {
	res := []TableFieldFull{}
	for _, v := range field {
		res = append(res, ToTableFieldFull(v))
	}
	return res
}
