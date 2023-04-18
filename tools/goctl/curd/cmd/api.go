package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"
	apiParser "github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/curd/util"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
	"os"
)

var (
	Api = cobrax.NewCommand("api", cobrax.WithRunE(runApi))
)

func init() {
	Api.Flags().StringVar(&api, "api")
	Api.Flags().StringVar(&url, "url")
	Api.Flags().StringVar(&table, "table")
	Api.Flags().StringVarWithDefaultValue(&style, "style", config.DefaultFormat)
}

func runApi(cmd *cobra.Command, args []string) error {
	if len(api) == 0 {
		return errors.New("missing -api")
	}
	if len(url) == 0 {
		return errors.New("missing -url")
	}
	if len(table) == 0 {
		return errors.New("missing -table")
	}

	apiSpec, err := apiParser.Parse(api)
	if err != nil {
		return err
	}

	tableInfo, err := util.ParseTable(url, table)
	if err != nil {
		return err
	}
	cfg, err := config.NewConfig(style)
	if err != nil {
		return err
	}

	if apiSpec, err = util.ReplaceApi(apiSpec, cfg, tableInfo); err != nil {
		return err
	}

	logx.Must(os.WriteFile(api, []byte(util.ApiSpecToString(apiSpec)), 0666))
	return nil
}
