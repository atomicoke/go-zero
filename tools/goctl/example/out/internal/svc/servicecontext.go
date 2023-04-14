package svc

import (
	"dm-admin/api/admin/internal/config"
	"dm-admin/api/admin/internal/middleware"
	"dm-admin/api/admin/internal/model"
	"dm-admin/common/asynqmq"
	"dm-admin/common/flow"
	"dm-admin/common/mid"
	"dm-admin/common/thirdparty/pay"
	"dm-admin/common/utils/userflow"

	"dm.com/toolx/arr"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Version                        int
	Config                         config.Config
	Redis                          *redis.Redis
	PermMenuAuth                   rest.Middleware
	DeleteVerifyCode               rest.Middleware
	LoginLog                       rest.Middleware
	SysUserModel                   model.SysUserModel
	SysDictionaryModel             model.SysDictionaryModel
	SysLogModel                    model.SysLogModel
	SysMenuModel                   model.SysMenuModel
	SysRoleModel                   model.SysRoleModel
	UserModel                      model.ThinkUsersModel
	GrayUserModel                  model.GrayUsersModel
	GrayDeviceModel                model.GrayDeviceModel
	IpLibraryModel                 model.ThinkIpLibraryModel
	UsersModel                     model.ThinkUsersModel
	PartnerCityUsersModel          model.PartnerCityUsersModel
	OfflineMeetingModel            model.OfflineMeetingModel
	SysActionLogModel              model.SysActionLogModel
	ThinkConfigModel               model.ThinkConfigModel
	FilesModel                     model.FilesModel
	RechargeModel                  model.ThinkRechargeModel
	RechargeModelV2                model.ThinkRechargeV2Model
	RechargeOrderModel             model.ThinkRechargeOrderModel
	GoodsModel                     model.ThinkGoodsModel
	GoodsSkuModel                  model.ThinkGoodsSkuModel
	GoodsOrderModel                model.ThinkGoodsOrderModel
	GoodsOrderInfoModel            model.ThinkGoodsOrderInfoModel
	GoodsCategoryModel             model.ThinkGoodsCategoryModel
	GoodsSpecRelModel              model.ThinkGoodsSpecRelModel
	UserCompleteTaskModel          model.ThinkUserCompleteTaskModel
	PartnerContributionModel       model.PartnerContributionModel
	UsersInfoModel                 model.ThinkUsersInfoModel
	PhoneDataV2Model               model.ThinkPhoneDataV2Model
	ManagementModel                model.ThinkManagementModel
	AreasModel                     model.AreasModel
	PayConfigModel                 model.ThinkPayConfigModel
	UsersFlowModel                 model.ThinkUsersFlowModel
	SpecModel                      model.ThinkSpecModel
	SpecValueModel                 model.ThinkSpecValueModel
	registerCount                  *registerCount
	UsersFlowService               userflow.IUserFlow
	IntegralRecordModel            model.ThinkIntegralRecordModel
	SowingModel                    model.ThinkSowingModel
	UsersAdvertiseModel            model.ThinkUsersAdvertiseModel
	SeedRecordModel                model.ThinkSeedRecordModel
	FarmLogModel                   model.ThinkFarmLogModel
	AppThemeModel                  model.AppThemeModel
	UsersBanModel                  model.ThinkUsersBanModel
	AppThemeNavBarModel            model.AppThemeNavBarModel
	ThinkUsersNavv2Model           model.ThinkUsersNavv2Model
	AuthFaceModel                  model.ThinkAuthFaceModel
	AppIconsModel                  model.AppIconsModel
	GoodsImageModel                model.ThinkGoodsImageModel
	PartnerContributionReviewModel model.PartnerContributionReviewModel
	TicketOrderModel               model.TicketOrderModel
	PublicServeModel               model.PublicServeModel
	AppThemeCategoryModel          model.AppThemeCategoryModel
	PartnerActiveUsersModel        model.PartnerActiveUsersModel
	OfflineMeetingUserModel        model.OfflineMeetingUserModel
	UsersActiveModel               model.ThinkUsersActiveModel
	DayTaskNoteModel               model.DayTaskNoteModel
	DayTotalModel                  model.DayTotalModel
	Pay                            *pay.Pay
	dayTaskNote                    *arr.MapX[string, *model.DayTaskNote]
	AdvMerchantModel               model.ThinkAdvMerchantModel
	AppCategoryModel               model.ThinkAppCategoryModel
	AppHomeTitleModel              model.AppHomeTitleModel
	HotelOrderModel                model.HotelOrderModel
	TasksModel                     model.TasksModel
	PartnerActiveUsersV2Model      model.PartnerActiveUsersV2Model
	UserFlowV2                     *flow.UserFlow
	CollectModel                   model.CollectModel
	PartnerActiveUsersV3Model      model.PartnerActiveUsersV3Model
	AdvMerchantRecordModel         model.AdvMerchantRecordModel
	AdvMerchantUserRecordModel     model.AdvMerchantUserRecordModel
	DisasterModel                  model.ThinkDisasterModel
	SowingDisasterModel            model.SowingDisasterModel
	MchPayConfigModel              model.MchPayConfigModel
	CommunityModel                 model.ThinkCommunityModel
	CommunityCategoryModel         model.ThinkCommunityCategoryModel
	FarmMaterialsModel             model.FarmMaterialsModel
	FarmMaterialsOrderModel        model.FarmMaterialsOrderModel
	PartnerCityActiveUsersModel    model.PartnerCityActiveUsersModel
	RechargePlatModel              model.ThinkRechargePlatModel
	Queue                          *asynqmq.AsynqClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	redisClient := redis.MustNewRedis(redis.RedisConf{
		Host: c.Redis.Host,
		Type: c.Redis.Type,
		Pass: c.Redis.Pass,
	})

	var s = &ServiceContext{
		Config:                         c,
		Redis:                          redisClient,
		PermMenuAuth:                   middleware.NewPermMenuAuthMiddleware(redisClient).Handle,
		SysUserModel:                   model.NewSysUserModel(mysqlConn, c.Cache),
		SysDictionaryModel:             model.NewSysDictionaryModel(mysqlConn, c.Cache),
		SysLogModel:                    model.NewSysLogModel(mysqlConn, c.Cache),
		SysMenuModel:                   model.NewSysMenuModel(mysqlConn, c.Cache),
		SysRoleModel:                   model.NewSysRoleModel(mysqlConn, c.Cache),
		UserModel:                      model.NewThinkUsersModel(mysqlConn, c.Cache),
		GrayUserModel:                  model.NewGrayUsersModel(mysqlConn, c.Cache),
		IpLibraryModel:                 model.NewThinkIpLibraryModel(mysqlConn, c.Cache),
		UsersModel:                     model.NewThinkUsersModel(mysqlConn, c.Cache),
		OfflineMeetingModel:            model.NewOfflineMeetingModel(mysqlConn, c.Cache),
		SysActionLogModel:              model.NewSysActionLogModel(mysqlConn, c.Cache),
		ThinkConfigModel:               model.NewThinkConfigModel(mysqlConn, c.Cache),
		FilesModel:                     model.NewFilesModel(mysqlConn, c.Cache),
		RechargeModel:                  model.NewThinkRechargeModel(mysqlConn, c.Cache),
		RechargeModelV2:                model.NewThinkRechargeV2Model(mysqlConn, c.Cache),
		RechargeOrderModel:             model.NewThinkRechargeOrderModel(mysqlConn, c.Cache),
		GoodsModel:                     model.NewThinkGoodsModel(mysqlConn, c.Cache),
		GoodsSkuModel:                  model.NewThinkGoodsSkuModel(mysqlConn, c.Cache),
		GoodsOrderModel:                model.NewThinkGoodsOrderModel(mysqlConn, c.Cache),
		GoodsOrderInfoModel:            model.NewThinkGoodsOrderInfoModel(mysqlConn, c.Cache),
		GoodsCategoryModel:             model.NewThinkGoodsCategoryModel(mysqlConn, c.Cache),
		GoodsSpecRelModel:              model.NewThinkGoodsSpecRelModel(mysqlConn, c.Cache),
		UserCompleteTaskModel:          model.NewThinkUserCompleteTaskModel(mysqlConn, c.Cache),
		PartnerContributionModel:       model.NewPartnerContributionModel(mysqlConn, c.Cache),
		UsersInfoModel:                 model.NewThinkUsersInfoModel(mysqlConn, c.Cache),
		ManagementModel:                model.NewThinkManagementModel(mysqlConn, c.Cache),
		AreasModel:                     model.NewAreasModel(mysqlConn, c.Cache),
		PayConfigModel:                 model.NewThinkPayConfigModel(mysqlConn, c.Cache),
		UsersFlowModel:                 model.NewThinkUsersFlowModel(mysqlConn, c.Cache),
		SpecModel:                      model.NewThinkSpecModel(mysqlConn, c.Cache),
		SpecValueModel:                 model.NewThinkSpecValueModel(mysqlConn, c.Cache),
		IntegralRecordModel:            model.NewThinkIntegralRecordModel(mysqlConn, c.Cache),
		SowingModel:                    model.NewThinkSowingModel(mysqlConn, c.Cache),
		UsersAdvertiseModel:            model.NewThinkUsersAdvertiseModel(mysqlConn, c.Cache),
		SeedRecordModel:                model.NewThinkSeedRecordModel(mysqlConn, c.Cache),
		FarmLogModel:                   model.NewThinkFarmLogModel(mysqlConn, c.Cache),
		UsersBanModel:                  model.NewThinkUsersBanModel(mysqlConn, c.Cache),
		DeleteVerifyCode:               middleware.NewDeleteVerifyCodeMiddleware(redisClient).Handle,
		AppThemeModel:                  model.NewAppThemeModel(mysqlConn, c.Cache),
		AppThemeNavBarModel:            model.NewAppThemeNavBarModel(mysqlConn, c.Cache),
		ThinkUsersNavv2Model:           model.NewThinkUsersNavv2Model(mysqlConn, c.Cache),
		AuthFaceModel:                  model.NewThinkAuthFaceModel(mysqlConn, c.Cache),
		AppIconsModel:                  model.NewAppIconsModel(mysqlConn, c.Cache),
		GoodsImageModel:                model.NewThinkGoodsImageModel(mysqlConn, c.Cache),
		PartnerContributionReviewModel: model.NewPartnerContributionReviewModel(mysqlConn, c.Cache),
		TicketOrderModel:               model.NewTicketOrderModel(mysqlConn, c.Cache),
		PublicServeModel:               model.NewPublicServeModel(mysqlConn, c.Cache),
		AppThemeCategoryModel:          model.NewAppThemeCategoryModel(mysqlConn, c.Cache),
		PartnerActiveUsersModel:        model.NewPartnerActiveUsersModel(mysqlConn, c.Cache),
		OfflineMeetingUserModel:        model.NewOfflineMeetingUserModel(mysqlConn, c.Cache),
		UsersActiveModel:               model.NewThinkUsersActiveModel(mysqlConn, c.Cache),
		GrayDeviceModel:                model.NewGrayDeviceModel(mysqlConn, c.Cache),
		PartnerCityUsersModel:          model.NewPartnerCityUsersModel(mysqlConn, c.Cache),
		DayTotalModel:                  model.NewDayTotalModel(mysqlConn, c.Cache),
		dayTaskNote:                    arr.NewMapX[string, *model.DayTaskNote](),
		DayTaskNoteModel:               model.NewDayTaskNoteModel(mysqlConn, c.Cache),
		AdvMerchantModel:               model.NewThinkAdvMerchantModel(mysqlConn, c.Cache),
		AppCategoryModel:               model.NewThinkAppCategoryModel(mysqlConn, c.Cache),
		AppHomeTitleModel:              model.NewAppHomeTitleModel(mysqlConn, c.Cache),
		HotelOrderModel:                model.NewHotelOrderModel(mysqlConn, c.Cache),
		TasksModel:                     model.NewTasksModel(mysqlConn, c.Cache),
		PartnerActiveUsersV2Model:      model.NewPartnerActiveUsersV2Model(mysqlConn, c.Cache),
		UserFlowV2:                     flow.NewUserFlow(mysqlConn, c.Cache),
		CollectModel:                   model.NewCollectModel(mysqlConn, c.Cache),
		PartnerActiveUsersV3Model:      model.NewPartnerActiveUsersV3Model(mysqlConn, c.Cache),
		AdvMerchantRecordModel:         model.NewAdvMerchantRecordModel(mysqlConn, c.Cache),
		AdvMerchantUserRecordModel:     model.NewAdvMerchantUserRecordModel(mysqlConn, c.Cache),
		DisasterModel:                  model.NewThinkDisasterModel(mysqlConn, c.Cache),
		SowingDisasterModel:            model.NewSowingDisasterModel(mysqlConn, c.Cache),
		MchPayConfigModel:              model.NewMchPayConfigModel(mysqlConn, c.Cache),
		FarmLiveModel:                  model.NewFarmLiveModel(mysqlConn, c.Cache),
		CommunityModel:                 model.NewThinkCommunityModel(mysqlConn, c.Cache),
		CommunityCategoryModel:         model.NewThinkCommunityCategoryModel(mysqlConn, c.Cache),
		FarmMaterialsModel:             model.NewFarmMaterialsModel(mysqlConn, c.Cache),
		FarmMaterialsOrderModel:        model.NewFarmMaterialsOrderModel(mysqlConn, c.Cache),
		PartnerCityActiveUsersModel:    model.NewPartnerCityActiveUsersModel(mysqlConn, c.Cache),
		PhoneDataV2Model:               model.NewThinkPhoneDataV2Model(mysqlConn, c.Cache),
		RechargePlatModel:              model.NewThinkRechargePlatModel(mysqlConn, c.Cache),
	}

	s.registerCount = NewRegisterCount(s)
	s.UsersFlowService = userflow.NewUserFlow(mysqlConn, c.Cache)
	s.LoginLog = mid.HandleProxy(middleware.NewLoginLogMiddleware(s.SysLogModel).Handle)
	s.Pay = pay.NewPayment()
	s.Queue = asynqmq.NewAsynqClient(c.Redis.Host, c.Redis.Pass)
	return s
}

func (s *ServiceContext) IsTest() bool {
	return s.Config.Mode == "test"
}
