package util

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/model"
	"strings"
)

func ParseTable(url string, table string) (*model.Table, error) {
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
