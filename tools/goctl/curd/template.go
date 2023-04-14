package curd

import "github.com/zeromicro/go-zero/tools/goctl/util/pathx"
import (
	_ "embed"
)

const (
	category        = "curd"
	addLogicFile    = "addLogic.tpl"
	deleteLogicFile = "deleteLogic.tpl"
	getLogicFile    = "getLogic.tpl"
	pageLogicFile   = "pageLogic.tpl"
	updateLogicFile = "updateLogic.tpl"
)

var (
	//go:embed tpl/addLogic.tpl
	addLogic string
	//go:embed tpl/deleteLogic.tpl
	deleteLogic string
	//go:embed tpl/getLogic.tpl
	getLogic string
	//go:embed tpl/pageLogic.tpl
	pageLogic string
	//go:embed tpl/updateLogic.tpl
	updateLogic string
)

var templates = map[string]string{
	addLogicFile:    addLogic,
	deleteLogicFile: deleteLogic,
	getLogicFile:    getLogic,
	pageLogicFile:   pageLogic,
	updateLogicFile: updateLogic,
}

func actionToLogicFile(s string) string {
	a := action(s)
	switch a {
	case Add:
		return addLogicFile
	case Delete:
		return deleteLogicFile
	case Get:
		return getLogicFile
	case Page:
		return pageLogicFile
	case Update:
		return updateLogicFile
	}

	return ""
}

// Category returns the category of the api files.
func Category() string {
	return category
}

// GenTemplates generates api template files.
func GenTemplates() error {
	return pathx.InitTemplates(category, templates)
}
