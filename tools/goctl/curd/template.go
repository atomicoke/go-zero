package curd

import "github.com/zeromicro/go-zero/tools/goctl/util/pathx"
import (
	_ "embed"
)

const (
	category = "curd"
)

//go:embed api.tpl
var apiTpl string

var templates = map[string]string{
	"api.tpl": apiTpl,
}

// Category returns the category of the api files.
func Category() string {
	return category
}

// GenTemplates generates api template files.
func GenTemplates() error {
	return pathx.InitTemplates(category, templates)
}
