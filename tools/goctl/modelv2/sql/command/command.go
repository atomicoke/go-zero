package command

import (
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/types"
	"path/filepath"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/modelv2/sql/command/migrationnotes"
	"github.com/zeromicro/go-zero/tools/goctl/modelv2/sql/gen"
	"github.com/zeromicro/go-zero/tools/goctl/modelv2/sql/model"
	"github.com/zeromicro/go-zero/tools/goctl/modelv2/sql/util"
	file "github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	// VarStringSrc describes the source file of sql.
	VarStringSrc string
	// VarStringDir describes the output directory of sql.
	VarStringDir string
	// VarBoolCache describes whether the cache is enabled.
	VarBoolCache bool
	// VarBoolIdea describes whether is idea or not.
	VarBoolIdea bool
	// VarStringURL describes the dsn of the sql.
	VarStringURL string
	// VarStringSliceTable describes tables.
	VarStringSliceTable []string
	// VarStringTable describes a table of sql.
	VarStringTable string
	// VarStringStyle describes the style.
	VarStringStyle string
	// VarStringDatabase describes the database.
	VarStringDatabase string
	// VarStringSchema describes the schema of postgresql.
	VarStringSchema string
	// VarStringHome describes the goctl home.
	VarStringHome string
	// VarStringRemote describes the remote git repository.
	VarStringRemote string
	// VarStringBranch describes the git branch of the repository.
	VarStringBranch string
	// VarBoolStrict describes whether the strict mode is enabled.
	VarBoolStrict bool
	// VarStringSliceIgnoreColumns represents the columns which are ignored.
	VarStringSliceIgnoreColumns []string
	// serviceContext 文件的目录绝对路径
	VarStringServiceContextDir string
	// new model 的参数
	VarStringModelArg string
)

var errNotMatched = errors.New("sql not matched")

// MysqlDDL generates modelv2 code from ddl
func MysqlDDL(_ *cobra.Command, _ []string) error {
	migrationnotes.BeforeCommands(VarStringDir, VarStringStyle)
	src := VarStringSrc
	dir := VarStringDir
	cache := VarBoolCache
	idea := VarBoolIdea
	style := VarStringStyle
	database := VarStringDatabase
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	if len(remote) > 0 {
		repo, _ := file.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}
	cfg, err := config.NewConfig(style)
	if err != nil {
		return err
	}

	if VarStringServiceContextDir == "" {
		return errors.New("service context dir is empty")
	}

	arg := types.DataSourceArg{
		SvcDir:        VarStringServiceContextDir,
		Src:           src,
		Dir:           dir,
		Cfg:           cfg,
		Cache:         cache,
		Idea:          idea,
		Database:      database,
		Strict:        VarBoolStrict,
		IgnoreColumns: mergeColumns(VarStringSliceIgnoreColumns),
		ModelArg:      VarStringModelArg,
	}
	return fromDDL(arg)
}

// MySqlDataSource generates modelv2 code from datasource
func MySqlDataSource(_ *cobra.Command, _ []string) error {
	migrationnotes.BeforeCommands(VarStringDir, VarStringStyle)
	url := strings.TrimSpace(VarStringURL)
	dir := strings.TrimSpace(VarStringDir)
	cache := VarBoolCache
	idea := VarBoolIdea
	style := VarStringStyle
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	if len(remote) > 0 {
		repo, _ := file.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	tableValue := VarStringSliceTable
	patterns := parseTableList(tableValue)
	cfg, err := config.NewConfig(style)
	if err != nil {
		return err
	}

	if VarStringServiceContextDir == "" {
		return errors.New("service context dir is empty")
	}

	arg := types.DataSourceArg{
		SvcDir:        VarStringServiceContextDir,
		Url:           url,
		Dir:           dir,
		TablePat:      patterns,
		Cfg:           cfg,
		Cache:         cache,
		Idea:          idea,
		Strict:        VarBoolStrict,
		IgnoreColumns: mergeColumns(VarStringSliceIgnoreColumns),
		ModelArg:      VarStringModelArg,
	}
	return fromMysqlDataSource(arg)
}

func mergeColumns(columns []string) []string {
	set := collection.NewSet()
	for _, v := range columns {
		fields := strings.FieldsFunc(v, func(r rune) bool {
			return r == ','
		})
		set.AddStr(fields...)
	}
	return set.KeysStr()
}

func parseTableList(tableValue []string) types.Pattern {
	tablePattern := make(types.Pattern)
	for _, v := range tableValue {
		fields := strings.FieldsFunc(v, func(r rune) bool {
			return r == ','
		})
		for _, f := range fields {
			tablePattern[f] = struct{}{}
		}
	}
	return tablePattern
}

// PostgreSqlDataSource generates modelv2 code from datasource
func PostgreSqlDataSource(_ *cobra.Command, _ []string) error {
	migrationnotes.BeforeCommands(VarStringDir, VarStringStyle)
	url := strings.TrimSpace(VarStringURL)
	dir := strings.TrimSpace(VarStringDir)
	cache := VarBoolCache
	idea := VarBoolIdea
	style := VarStringStyle
	schema := VarStringSchema
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	if len(remote) > 0 {
		repo, _ := file.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	if len(schema) == 0 {
		schema = "public"
	}

	pattern := strings.TrimSpace(VarStringTable)
	cfg, err := config.NewConfig(style)
	if err != nil {
		return err
	}

	arg := types.DataSourceArg{
		Src:           url,
		Dir:           dir,
		Cfg:           cfg,
		Cache:         cache,
		Idea:          idea,
		Database:      schema,
		Strict:        VarBoolStrict,
		IgnoreColumns: mergeColumns(VarStringSliceIgnoreColumns),
		SvcDir:        VarStringServiceContextDir,
	}

	return fromPostgreSqlDataSource(arg, pattern)
}

func fromDDL(arg types.DataSourceArg) error {
	log := console.NewConsole(arg.Idea)
	src := strings.TrimSpace(arg.Src)
	if len(src) == 0 {
		return errors.New("expected path or path globbing patterns, but nothing found")
	}

	files, err := util.MatchFiles(src)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errNotMatched
	}

	generator, err := gen.NewDefaultGenerator(arg.Dir, arg.Cfg,
		gen.WithConsoleOption(log), gen.WithIgnoreColumns(arg.IgnoreColumns))
	if err != nil {
		return err
	}

	for _, file := range files {
		err = generator.StartFromDDL(file, arg.Cache, arg.Strict, arg.Database)
		if err != nil {
			return err
		}
	}

	return nil
}

type dataSourceArg struct {
	src           string
	url, dir      string
	tablePat      types.Pattern
	cfg           *config.Config
	cache, idea   bool
	strict        bool
	ignoreColumns []string
	svcDir        string
	database      string
}

func fromMysqlDataSource(arg types.DataSourceArg) error {
	log := console.NewConsole(arg.Idea)
	if len(arg.Url) == 0 {
		log.Error("%v", "expected data source of mysql, but nothing found")
		return nil
	}

	if len(arg.TablePat) == 0 {
		log.Error("%v", "expected table or table globbing patterns, but nothing found")
		return nil
	}

	dsn, err := mysql.ParseDSN(arg.Url)
	if err != nil {
		return err
	}

	logx.Disable()
	databaseSource := strings.TrimSuffix(arg.Url, "/"+dsn.DBName) + "/information_schema"
	db := sqlx.NewMysql(databaseSource)
	im := model.NewInformationSchemaModel(db)

	tables, err := im.GetAllTables(dsn.DBName)
	if err != nil {
		return err
	}

	matchTables := make(map[string]*model.Table)
	for _, item := range tables {
		if !arg.TablePat.Match(item) {
			continue
		}

		columnData, err := im.FindColumns(dsn.DBName, item)
		if err != nil {
			return err
		}

		table, err := columnData.Convert()
		if err != nil {
			return err
		}

		matchTables[item] = table
	}

	if len(matchTables) == 0 {
		return errors.New("no tables matched")
	}

	generator, err := gen.NewDefaultGenerator(arg.Dir, arg.Cfg,
		gen.WithConsoleOption(log), gen.WithIgnoreColumns(arg.IgnoreColumns))
	if err != nil {
		return err
	}

	return generator.StartFromInformationSchema(matchTables, arg)
}

func fromPostgreSqlDataSource(arg types.DataSourceArg, pattern string) error {
	url := strings.TrimSpace(arg.Src)
	dir := strings.TrimSpace(arg.Dir)
	schema := strings.TrimSpace(arg.Database)
	cfg := arg.Cfg
	idea := arg.Idea

	log := console.NewConsole(idea)
	if len(url) == 0 {
		log.Error("%v", "expected data source of postgresql, but nothing found")
		return nil
	}

	if len(pattern) == 0 {
		log.Error("%v", "expected table or table globbing patterns, but nothing found")
		return nil
	}
	db := postgres.New(url)
	im := model.NewPostgreSqlModel(db)

	tables, err := im.GetAllTables(schema)
	if err != nil {
		return err
	}

	matchTables := make(map[string]*model.Table)
	for _, item := range tables {
		match, err := filepath.Match(pattern, item)
		if err != nil {
			return err
		}

		if !match {
			continue
		}

		columnData, err := im.FindColumns(schema, item)
		if err != nil {
			return err
		}

		table, err := columnData.Convert()
		if err != nil {
			return err
		}

		matchTables[item] = table
	}

	if len(matchTables) == 0 {
		return errors.New("no tables matched")
	}

	generator, err := gen.NewDefaultGenerator(dir, cfg, gen.WithConsoleOption(log), gen.WithPostgreSql())
	if err != nil {
		return err
	}

	return generator.StartFromInformationSchema(matchTables, arg)
}
