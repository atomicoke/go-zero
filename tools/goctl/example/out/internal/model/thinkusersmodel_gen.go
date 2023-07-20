// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	thinkUsersFieldNames          = builder.RawFieldNames(&ThinkUsers{})
	thinkUsersRows                = strings.Join(thinkUsersFieldNames, ",")
	thinkUsersRowsExpectAutoSet   = strings.Join(stringx.Remove(thinkUsersFieldNames, "`user_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	thinkUsersRowsWithPlaceHolder = strings.Join(stringx.Remove(thinkUsersFieldNames, "`user_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	thinkUsersModel interface {
		Insert(ctx context.Context, data *ThinkUsers) (sql.Result, error)
		FindOne(ctx context.Context, userId int64) (*ThinkUsers, error)
		FindOneByPhone(ctx context.Context, phone string) (*ThinkUsers, error)
		Update(ctx context.Context, data *ThinkUsers) error
		Delete(ctx context.Context, userId int64) error
	}

	defaultThinkUsersModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ThinkUsers struct {
		UserId                   int64          `db:"user_id"`
		ParentId                 int64          `db:"parent_id"`                   // 上级user_id
		UserCode                 string         `db:"user_code"`                   // 用户编号
		Token                    string         `db:"token"`                       // 用户token
		StarLevel                int64          `db:"star_level"`                  // 星级
		Password                 string         `db:"password"`                    // 密码
		PayPassword              sql.NullString `db:"pay_password"`                // 支付密码
		Code                     string         `db:"code"`                        // 身份证
		Wechat                   sql.NullString `db:"wechat"`                      // 微信号
		NickName                 sql.NullString `db:"nick_name"`                   // 用户昵称
		RealName                 sql.NullString `db:"real_name"`                   // 真实姓名
		Portrait                 sql.NullString `db:"portrait"`                    // 头像
		Gender                   int64          `db:"gender"`                      // 性别：0未知、1男、2女
		Phone                    string         `db:"phone"`                       // 手机号码
		Integral                 float64        `db:"integral"`                    // 账户积分
		LibraryIntegral          float64        `db:"library_integral"`            // 仓库积分
		LibraryGold              float64        `db:"library_gold"`                // 仓库金豆
		LibrarySilver            float64        `db:"library_silver"`              // 仓库银豆
		Gold                     float64        `db:"gold"`                        // 农场金豆
		Frozen                   float64        `db:"frozen"`                      // 冻结金额
		Silver                   float64        `db:"silver"`                      // 农场银豆
		GoldPreProfit            float64        `db:"gold_pre_profit"`             // 金豆预估收益
		ExtractableGoldPreProfit float64        `db:"extractable_gold_pre_profit"` // 可提取的预估收益
		SilverPreProfit          float64        `db:"silver_pre_profit"`           // 银豆预估收益
		Status                   int64          `db:"status"`                      // 状态：1正常、2锁定
		Active                   float64        `db:"active"`                      // 活跃值
		RecommendActive          float64        `db:"recommend_active"`            // 推荐活跃值
		TeamActive               int64          `db:"team_active"`                 // 团队总活跃人数
		TeamNum                  int64          `db:"team_num"`                    // 团队总人数
		TeamAuth                 int64          `db:"team_auth"`                   // 团队认证人数
		Parents                  string         `db:"parents"`                     // 父祖关系
		YesterdayTeam            float64        `db:"yesterday_team"`              // 昨日团队贡献
		YesterdayOneTeam         float64        `db:"yesterday_one_team"`          // 1团昨日贡献值
		YesterdayTwoTeam         float64        `db:"yesterday_two_team"`          // 2团昨日贡献值
		Contribution             float64        `db:"contribution"`                // 领取任务贡献值总和与团队贡献之和
		TeamContribution         float64        `db:"team_contribution"`           // 团队贡献值
		IsTask                   int64          `db:"is_task"`                     // 用户是否有任务：1是、2否
		TaskNum                  int64          `db:"task_num"`                    // 一共兑换了多少任务
		IsParent                 int64          `db:"is_parent"`                   // 用户是否有上级：1是、2否
		OneNum                   int64          `db:"one_num"`                     // 1级人数
		TwoNum                   int64          `db:"two_num"`                     // 2级人数
		IsAuth                   int64          `db:"is_auth"`                     // 是否认证：1已认证、2未认证
		IsCheck                  int64          `db:"is_check"`                    // 是否需要核验：1是、2否
		IsShop                   int64          `db:"is_shop"`                     // 是否为商家：1是、2否
		CheckType                int64          `db:"check_type"`                  // 验证类型：1短信
		Model                    sql.NullString `db:"model"`                       // 设备信息
		PassTime                 int64          `db:"pass_time"`                   // 修改密码锁定时间
		CreateTime               int64          `db:"create_time"`                 // 创建时间
		UpdateTime               int64          `db:"update_time"`                 // 更新时间
		NumberStarts             int64          `db:"number_starts"`               // 启动次数
		Version                  sql.NullString `db:"version"`                     // 客户端使用的版本号
		LastLogin                sql.NullInt64  `db:"last_login"`                  // 最后操作时间
		TransferSms              int64          `db:"transfer_sms"`                //  是否需要转赠短信验证：1需要、2不需要
		SowingFirstTime          int64          `db:"sowing_first_time"`           // 首次播种时间
		SelfContribution         float64        `db:"self_contribution"`           // 自身贡献
		GetContribution          float64        `db:"get_contribution"`            // 获得贡献
		TotalContribution        float64        `db:"total_contribution"`          // 总贡献
		AuthAt                   time.Time      `db:"auth_at"`                     // 实人认证通过时间
		LockIntegral             int64          `db:"lock_integral"`
		ArtificialAuth           int64          `db:"artificial_auth"` // 1:人为授权
		PhoneCheck               int64          `db:"phone_check"`     // 手机验证，1：已验证，0未验证
		Udid                     string         `db:"udid"`            // 设备唯一码
		V                        int64          `db:"v"`               // 版本锁
		LandRemain               float64        `db:"land_remain"`     // 土地剩余收获
		BondAmount               float64        `db:"bond_amount"`     // 信誉度
	}
)

func newThinkUsersModel(conn sqlx.SqlConn) *defaultThinkUsersModel {
	return &defaultThinkUsersModel{
		conn:  conn,
		table: "`think_users`",
	}
}

func (m *defaultThinkUsersModel) Delete(ctx context.Context, userId int64) error {
	query := fmt.Sprintf("delete from %s where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userId)
	return err
}

func (m *defaultThinkUsersModel) FindOne(ctx context.Context, userId int64) (*ThinkUsers, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", thinkUsersRows, m.table)
	var resp ThinkUsers
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultThinkUsersModel) FindOneByPhone(ctx context.Context, phone string) (*ThinkUsers, error) {
	var resp ThinkUsers
	query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", thinkUsersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, phone)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultThinkUsersModel) Insert(ctx context.Context, data *ThinkUsers) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, thinkUsersRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.ParentId, data.UserCode, data.Token, data.StarLevel, data.Password, data.PayPassword, data.Code, data.Wechat, data.NickName, data.RealName, data.Portrait, data.Gender, data.Phone, data.Integral, data.LibraryIntegral, data.LibraryGold, data.LibrarySilver, data.Gold, data.Frozen, data.Silver, data.GoldPreProfit, data.ExtractableGoldPreProfit, data.SilverPreProfit, data.Status, data.Active, data.RecommendActive, data.TeamActive, data.TeamNum, data.TeamAuth, data.Parents, data.YesterdayTeam, data.YesterdayOneTeam, data.YesterdayTwoTeam, data.Contribution, data.TeamContribution, data.IsTask, data.TaskNum, data.IsParent, data.OneNum, data.TwoNum, data.IsAuth, data.IsCheck, data.IsShop, data.CheckType, data.Model, data.PassTime, data.NumberStarts, data.Version, data.LastLogin, data.TransferSms, data.SowingFirstTime, data.SelfContribution, data.GetContribution, data.TotalContribution, data.AuthAt, data.LockIntegral, data.ArtificialAuth, data.PhoneCheck, data.Udid, data.V, data.LandRemain, data.BondAmount)
	return ret, err
}

func (m *defaultThinkUsersModel) Update(ctx context.Context, newData *ThinkUsers) error {
	query := fmt.Sprintf("update %s set %s where `user_id` = ?", m.table, thinkUsersRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.ParentId, newData.UserCode, newData.Token, newData.StarLevel, newData.Password, newData.PayPassword, newData.Code, newData.Wechat, newData.NickName, newData.RealName, newData.Portrait, newData.Gender, newData.Phone, newData.Integral, newData.LibraryIntegral, newData.LibraryGold, newData.LibrarySilver, newData.Gold, newData.Frozen, newData.Silver, newData.GoldPreProfit, newData.ExtractableGoldPreProfit, newData.SilverPreProfit, newData.Status, newData.Active, newData.RecommendActive, newData.TeamActive, newData.TeamNum, newData.TeamAuth, newData.Parents, newData.YesterdayTeam, newData.YesterdayOneTeam, newData.YesterdayTwoTeam, newData.Contribution, newData.TeamContribution, newData.IsTask, newData.TaskNum, newData.IsParent, newData.OneNum, newData.TwoNum, newData.IsAuth, newData.IsCheck, newData.IsShop, newData.CheckType, newData.Model, newData.PassTime, newData.NumberStarts, newData.Version, newData.LastLogin, newData.TransferSms, newData.SowingFirstTime, newData.SelfContribution, newData.GetContribution, newData.TotalContribution, newData.AuthAt, newData.LockIntegral, newData.ArtificialAuth, newData.PhoneCheck, newData.Udid, newData.V, newData.LandRemain, newData.BondAmount, newData.UserId)
	return err
}

func (m *defaultThinkUsersModel) tableName() string {
	return m.table
}