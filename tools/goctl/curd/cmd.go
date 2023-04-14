package curd

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	apiParser "github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/model"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"strings"
)

var (
	Cmd = cobrax.NewCommand("curd", cobrax.WithRunE(runE))
	// The api file
	api string
	// The target dir
	dir string
	// The table name
	table string
	// The data source of database,like "root:password@tcp(127.0.0.1:3306)/database
	url string
	// The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority
	home string
	// describes the style of output files.
	style string
)

func init() {
	Cmd.Flags().StringVar(&api, "api")
	Cmd.Flags().StringVar(&dir, "dir")
	Cmd.Flags().StringVar(&url, "url")
	Cmd.Flags().StringVar(&table, "table")
	Cmd.Flags().StringVar(&home, "home")
	Cmd.Flags().StringVarWithDefaultValue(&style, "style", config.DefaultFormat)
}

func runE(cmd *cobra.Command, args []string) error {
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}
	if len(api) == 0 {
		return errors.New("missing -api")
	}
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}
	if len(url) == 0 {
		return errors.New("missing -url")
	}
	if len(table) == 0 {
		return errors.New("missing -table")
	}
	return doGenCrud(api, dir, url, table, style)
}

// doGenCrud gen crud files with api file
func doGenCrud(apiFile, dir, url, table, namingStyle string) error {
	apiSpec, err := apiParser.Parse(apiFile)
	if err != nil {
		return err
	}

	tableInfo, err := parseTable(url, table)
	if err != nil {
		return err
	}

	cfg, err := config.NewConfig(namingStyle)
	if err != nil {
		return err
	}
	logx.Must(pathx.MkdirIfNotExist(dir))
	rootPkg, err := golang.GetParentPackage(dir)
	if err != nil {
		return err
	}

	if apiSpec, err = replaceApi(apiSpec, cfg, tableInfo); err != nil {
		return err
	}

	logx.Must(genTypes(dir, cfg, apiSpec))
	logx.Must(gogen.GenHandlers(dir, rootPkg, cfg, apiSpec))
	logx.Must(genLogic(dir, rootPkg, cfg, apiSpec))

	//fmt.Println(apiSpecToString(apiSpec))
	return nil
}

func parseTable(url string, table string) (*model.Table, error) {
	dsn, err := mysql.ParseDSN(url)
	if err != nil {
		return nil, err
	}
	databaseSource := strings.TrimSuffix(url, "/"+dsn.DBName) + "/information_schema"
	db := sqlx.NewMysql(databaseSource)
	im := model.NewInformationSchemaModel(db)
	tables, err := im.GetAllTables(dsn.DBName)
	if err != nil {
		return nil, err
	}

	for _, item := range tables {
		if table == item {
			columnData, err := im.FindColumns(dsn.DBName, item)
			if err != nil {
				return nil, err
			}

			return columnData.Convert()
		}
	}
	return nil, errors.New("table not found")
}
