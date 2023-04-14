package types

import (
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"path/filepath"
)

type Pattern map[string]struct{}

type DataSourceArg struct {
	Src           string
	Url, Dir      string
	TablePat      Pattern
	Cfg           *config.Config
	Cache, Idea   bool
	Strict        bool
	IgnoreColumns []string
	SvcDir        string
	Database      string
	ModelArg      string
}

func (p Pattern) Match(s string) bool {
	for v := range p {
		match, err := filepath.Match(v, s)
		if err != nil {
			console.Error("%+v", err)
			continue
		}
		if match {
			return true
		}
	}
	return false
}

func (p Pattern) List() []string {
	var ret []string
	for v := range p {
		ret = append(ret, v)
	}
	return ret
}
