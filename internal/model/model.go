package model

import (
	"WowjoyProject/ObjectCloudService/pkg/setting"
	"database/sql"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type KeyData struct {
	instance_key, img_upload_status, dcm_upload_status          sql.NullInt64
	imgfile, remoteimgfile, dcmfile, remotedcmfile, ip, virpath sql.NullString
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*sql.DB, error) {
	db, err := sql.Open(databaseSetting.DBType, databaseSetting.DBConn)
	if err != nil {
		return nil, err
	}
	// 数据库最大连接数
	db.SetMaxOpenConns(databaseSetting.MaxIdleConns)
	db.SetMaxIdleConns(databaseSetting.MaxIdleConns)

	return db, nil
}
