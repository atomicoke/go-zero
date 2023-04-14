package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ThinkUsersModel = (*customThinkUsersModel)(nil)

type (
	// ThinkUsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customThinkUsersModel.
	ThinkUsersModel interface {
		thinkUsersModel
	}

	customThinkUsersModel struct {
		*defaultThinkUsersModel
	}
)

// NewThinkUsersModel returns a model for the database table.
func NewThinkUsersModel(conn sqlx.SqlConn, cache interface{}) ThinkUsersModel {
	return &customThinkUsersModel{
		defaultThinkUsersModel: newThinkUsersModel(conn),
	}
}
