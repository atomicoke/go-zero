package tpl

import "github.com/zeromicro/go-zero/tools/goctl/util/pathx"
import (
	_ "embed"
)

const (
	category        = "curd"
	AddLogicFile    = "addLogic.tpl"
	DeleteLogicFile = "deleteLogic.tpl"
	GetLogicFile    = "getLogic.tpl"
	PageLogicFile   = "pageLogic.tpl"
	UpdateLogicFile = "updateLogic.tpl"
)

var (
	//go:embed addLogic.tpl
	addLogic string
	//go:embed deleteLogic.tpl
	deleteLogic string
	//go:embed getLogic.tpl
	getLogic string
	//go:embed pageLogic.tpl
	pageLogic string
	//go:embed updateLogic.tpl
	updateLogic string
)

var Templates = map[string]string{
	AddLogicFile:    addLogic,
	DeleteLogicFile: deleteLogic,
	GetLogicFile:    getLogic,
	PageLogicFile:   pageLogic,
	UpdateLogicFile: updateLogic,
}

func ActionToLogicFile(s string) string {
	a := Action(s)
	switch a {
	case Add:
		return AddLogicFile
	case Delete:
		return DeleteLogicFile
	case Get:
		return GetLogicFile
	case Page:
		return PageLogicFile
	case Update:
		return UpdateLogicFile
	}

	return ""
}

// Category returns the category of the api files.
func Category() string {
	return category
}

// GenTemplates generates api template files.
func GenTemplates() error {
	return pathx.InitTemplates(category, Templates)
}

type Action string

const (
	Add    Action = "add"
	Update Action = "update"
	Delete Action = "delete"
	Page   Action = "page"
	Get    Action = "get"
)

var (
	ActionMap = map[string]bool{
		string(Add):    true,
		string(Update): true,
		string(Delete): true,
		string(Get):    true,
		string(Page):   true,
	}
)
