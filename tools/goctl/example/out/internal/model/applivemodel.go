package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AppLiveModel = (*customAppLiveModel)(nil)

type (
	// AppLiveModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAppLiveModel.
	AppLiveModel interface {
		appLiveModel
	}

	customAppLiveModel struct {
		*defaultAppLiveModel
	}
)

// NewAppLiveModel returns a model for the database table.
func NewAppLiveModel(conn sqlx.SqlConn, c interface{}) AppLiveModel {
	return &customAppLiveModel{
		defaultAppLiveModel: newAppLiveModel(conn),
	}
}
