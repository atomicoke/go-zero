package util

import (
	"fmt"
	"github.com/iancoleman/strcase"
)

func maptype(dataType string) string {
	switch dataType {
	case "time.Time", "sql.NullTime", "sql.NullString":
		return "string"
	case "sql.NullInt64", "sql.NullInt32", "sql.NullInt16", "sql.NullInt8", "uint64", "uint32", "uint16", "uint8":
		return "int64"
	case "sql.NullFloat64":
		return "float64"
	case "sql.NullBool":
		return "bool"
	default:
		return dataType
	}
}

var mapJsonTag = func(name string, comment string) string {
	return fmt.Sprintf("`label:\"%s\" json:\"%s\"`", comment, strcase.ToLowerCamel(name))
}

var mapFormTag = func(name string, comment string) string {
	return fmt.Sprintf("`label:\"%s\" form:\"%s\"`", comment, strcase.ToLowerCamel(name))
}
var mapFormTagWithValid = func(name string, comment string, valid string) string {
	return fmt.Sprintf("`label:\"%s\" validate:\"%s\" form:\"%s\"`", comment, valid, strcase.ToLowerCamel(name))
}
