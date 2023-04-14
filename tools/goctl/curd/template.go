package curd

import "github.com/zeromicro/go-zero/tools/goctl/util/pathx"
import (
	_ "embed"
)

const (
	category = "curd"
)

var (
	//go:embed tpl/add-logic.tpl
	addLogic string
	//go:embed tpl/delete-logic.tpl
	deleteLogic string
	//go:embed tpl/get-logic.tpl
	getLogic string
	//go:embed tpl/page-logic.tpl
	pageLogic string
	//go:embed tpl/update-logic.tpl
	updateLogic string
)

var templates = map[string]string{}

// Category returns the category of the api files.
func Category() string {
	return category
}

// GenTemplates generates api template files.
func GenTemplates() error {
	return pathx.InitTemplates(category, templates)
}
