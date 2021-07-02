package global

import (
	"WowjoyProject/ObjectCloudService/pkg/logger"
	"WowjoyProject/ObjectCloudService/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	GeneralSetting  *setting.GeneralSettingS
	DatabaseSetting *setting.DatabaseSettingS
	ObjectSetting   *setting.ObjectSettingS
	Logger          *logger.Logger
)
