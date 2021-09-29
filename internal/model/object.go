package model

import (
	"WowjoyProject/ObjectCloudService/global"
	"WowjoyProject/ObjectCloudService/pkg/general"
	"strings"
)

// test upload
func TestUploadeData(action global.ActionType) {
	data := global.ObjectData{
		InstanceKey: 123456,
		Key:         "pacs/ct56/2020/05/06/CT/CT.ed6b4e41c25c624f14f85410ba506afe.dcm",
		Type:        action,
		Path:        "W:\\image\\1234567.dcm",
	}
	global.ObjectDataChan <- data
}

// test Down
func TestDownData(action global.ActionType) {
	data := global.ObjectData{
		InstanceKey: 123456,
		Key:         "pacs/ct56/2020/05/06/CT/CT.ed6b4e41c25c624f14f85410ba506afe.dcm",
		Type:        action,
		Path:        "D:\\image\\Down\\CT.ed6b4e41c25c624f14f85410ba506afe.dcm",
	}
	global.ObjectDataChan <- data
}

// 通过检查号获取数据
func GetObjectData(accessNumber string, action global.ActionType) {
	sql := `select im.instance_key,im.img_file_name,im.img_file_name_remote,im.dcm_file_name_remote,ins.file_name,stu.ip,stu.s_virtual_dir
	from  image im 
	inner join instance ins on im.instance_key = ins.instance_key
	inner join study_location stu on ins.location_code = stu.n_station_code
	inner join study st on st.study_key = ins.study_key
	where st.accession_number=?;`
	// global.Logger.Debug(sql)
	rows, err := global.DBEngine.Query(sql, accessNumber)
	if err != nil {
		global.Logger.Fatal(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		key := KeyData{}
		_ = rows.Scan(&key.instance_key, &key.imgfile, &key.remoteimgfile, &key.remotedcmfile, &key.dcmfile, &key.ip, &key.virpath)
		if key.imgfile.String != "" {
			file_key, file_path := general.GetFilePath(action, key.imgfile.String, key.remoteimgfile.String, key.ip.String, key.virpath.String)
			global.Logger.Info(action, " 需要处理的文件名：", file_path)
			data := global.ObjectData{
				InstanceKey: key.instance_key.Int64,
				Key:         file_key,
				Type:        action,
				Path:        file_path,
				Count:       1,
			}
			global.ObjectDataChan <- data
		}
		if key.dcmfile.String != "" {
			file_key, file_path := general.GetFilePath(action, key.dcmfile.String, key.remotedcmfile.String, key.ip.String, key.virpath.String)
			global.Logger.Info(action, " 需要处理的文件名：", file_path)
			data := global.ObjectData{
				InstanceKey: key.instance_key.Int64,
				Key:         file_key,
				Type:        action,
				Path:        file_path,
				Count:       1,
			}
			global.ObjectDataChan <- data
		}
	}
}

// 上传数据后更新数据库
func UpdateUplaode(key int64, file string, value bool) {
	// 更新image表中 dcm_file_upload_status 和 dcm_file_name_remote
	// value true 上传成功
	// value false 上传失败
	success_code := global.ObjectSetting.OBJECT_Upload_SUCCESS
	error_code := global.ObjectSetting.OBJECT_Upload_ERROR
	upload_success_code := global.ObjectSetting.OBJECT_Upload_Success_Code
	if value {
		if strings.Contains(file, ".dcm") {
			sql := `update image im set im.dcm_file_upload_status=?,im.dcm_file_name_remote=? where im.instance_key=?`
			global.DBEngine.Exec(sql, success_code, file, key)
			// 更新上传成功的study_location
			sql = `update instance ins set ins.location_code=? where ins.instance_key=?`
			global.DBEngine.Exec(sql, upload_success_code, key)
		} else {
			sql := `update image im set im.img_file_upload_status=?,im.img_file_name_remote=? where im.instance_key=?`
			global.DBEngine.Exec(sql, success_code, file, key)
		}
	} else {
		if strings.Contains(file, ".dcm") {
			sql := `update image im set im.dcm_file_upload_status=? where im.instance_key=? and im.dcm_file_upload_status!=?`
			global.DBEngine.Exec(sql, error_code, key, success_code)
		} else {
			sql := `update image im set im.img_file_upload_status=? where im.instance_key=? and im.img_file_upload_status!=?`
			global.DBEngine.Exec(sql, error_code, key, success_code)
		}
	}
}

// 数据下载成功更新数据库
func UpdateDown(key int64, file string, value bool) {
	// 更新表中
	// value true 上传成功
	// value false 上传失败
	success_code := global.ObjectSetting.OBJECT_Down_SUCCESS
	error_code := global.ObjectSetting.OBJECT_Down_ERROR
	destcode := global.ObjectSetting.DOWN_Dest_Code
	if value {
		if strings.Contains(file, ".dcm") {
			sql := `update instance ins set ins.FileExist=?,ins.location_code=? where ins.instance_key=?`
			global.DBEngine.Exec(sql, success_code, destcode, key)
		}
	} else {
		if strings.Contains(file, ".dcm") {
			sql := `update instance ins set ins.FileExist=? where ins.instance_key=?`
			global.DBEngine.Exec(sql, error_code, key)
		}
	}
}

// 自动上传数据
func AutoUploadObjectData() {
	global.Logger.Info("******开始自动上传数据******")
	sql := `select im.instance_key,im.img_file_upload_status,im.dcm_file_upload_status,im.img_file_name,im.img_file_name_remote,im.dcm_file_name_remote,ins.file_name,stu.ip,stu.s_virtual_dir
		from  image im 
		inner join instance ins on im.instance_key = ins.instance_key
		inner join study_location stu on ins.location_code = stu.n_station_code
		where im.dcm_file_upload_status=? or im.img_file_upload_status=? limit ?;`
	// global.Logger.Debug(sql)
	rows, err := global.DBEngine.Query(sql, global.ObjectSetting.OBJECT_Upload_Flag, global.ObjectSetting.OBJECT_Upload_Flag, global.ObjectSetting.OBJECT_TASK)
	if err != nil {
		global.Logger.Fatal(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		key := KeyData{}
		_ = rows.Scan(&key.instance_key, &key.img_upload_status, &key.dcm_upload_status, &key.imgfile, &key.remoteimgfile, &key.remotedcmfile, &key.dcmfile, &key.ip, &key.virpath)
		if key.imgfile.String != "" && key.img_upload_status.Int64 == int64(global.ObjectSetting.OBJECT_Upload_Flag) {
			file_key, file_path := general.GetFilePath(global.UPLOAD, key.imgfile.String, key.remoteimgfile.String, key.ip.String, key.virpath.String)
			global.Logger.Info(global.UPLOAD, " 需要处理的文件名：", file_path)
			data := global.ObjectData{
				InstanceKey: key.instance_key.Int64,
				Key:         file_key,
				Type:        global.UPLOAD,
				Path:        file_path,
				Count:       1,
			}
			global.ObjectDataChan <- data
		}
		if key.imgfile.String == "" {
			global.Logger.Info("获取的jpg不存在，更新数据库....")
			UpdateEmptyPathJPG(key.instance_key.Int64)
		}
		if key.dcmfile.String != "" && key.dcm_upload_status.Int64 == int64(global.ObjectSetting.OBJECT_Upload_Flag) {
			file_key, file_path := general.GetFilePath(global.UPLOAD, key.dcmfile.String, key.remotedcmfile.String, key.ip.String, key.virpath.String)
			global.Logger.Info(global.UPLOAD, " 需要处理的文件名：", file_path)
			data := global.ObjectData{
				InstanceKey: key.instance_key.Int64,
				Key:         file_key,
				Type:        global.UPLOAD,
				Path:        file_path,
				Count:       1,
			}
			global.ObjectDataChan <- data
		}
	}
}

// 自动下载数据 // 文件上传成功，并且不存在
func AutoDownObjectData() {
	global.Logger.Info("******开始自动下载数据******")
	for global.AutoDowndFlag {
		sql := `select im.instance_key,im.img_file_name,im.img_file_name_remote,im.dcm_file_name_remote,ins.file_name
		from  image im 
		inner join instance ins on im.instance_key = ins.instance_key
		where im.dcm_file_upload_status= 0 and ins.FileExist=? limit ?;`
		// global.Logger.Debug(sql)
		rows, err := global.DBEngine.Query(sql, global.ObjectSetting.OBJECT_Down_Flag, global.ObjectSetting.OBJECT_TASK)
		if err != nil {
			global.Logger.Fatal(err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			key := KeyData{}
			_ = rows.Scan(&key.instance_key, &key.imgfile, &key.remoteimgfile, &key.remotedcmfile, &key.dcmfile)
			if key.imgfile.String != "" {
				file_key, file_path := general.GetFilePath(global.DOWNLOAD, key.imgfile.String, key.remoteimgfile.String, "", "")
				global.Logger.Info(global.DOWNLOAD, " 需要处理的文件名：", file_path)
				data := global.ObjectData{
					InstanceKey: key.instance_key.Int64,
					Key:         file_key,
					Type:        global.DOWNLOAD,
					Path:        file_path,
					Count:       1,
				}
				global.ObjectDataChan <- data
			}
			if key.dcmfile.String != "" {
				file_key, file_path := general.GetFilePath(global.DOWNLOAD, key.dcmfile.String, key.remotedcmfile.String, key.ip.String, key.virpath.String)
				global.Logger.Info(global.DOWNLOAD, " 需要处理的文件名：", file_path)
				data := global.ObjectData{
					InstanceKey: key.instance_key.Int64,
					Key:         file_key,
					Type:        global.DOWNLOAD,
					Path:        file_path,
					Count:       1,
				}
				global.ObjectDataChan <- data
			}
		}
	}
}

// 更新不存在的JPG字段
func UpdateEmptyPathJPG(key int64) {
	sql := `image im set im.img_file_upload_status = 2 where ins.instance_key=?`
	global.DBEngine.Exec(sql, key)
}
