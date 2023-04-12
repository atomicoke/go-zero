package main

import (
	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/cmd"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
)

var BuildTime = ""

func main() {
	version.BuildTime = BuildTime
	logx.Disable()
	load.Disable()
	cmd.Execute()
}
