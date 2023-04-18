package curd

import (
	"github.com/zeromicro/go-zero/tools/goctl/curd/cmd"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
)

var (
	Cmd = cobrax.NewCommand("curd")
)

func init() {
	Cmd.AddCommand(cmd.Api)
	Cmd.AddCommand(cmd.Gen)
}
