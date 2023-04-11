package object

import (
	"WowjoyProject/ObjectCloudService/global"
	"WowjoyProject/ObjectCloudService/internal/model"
	"WowjoyProject/ObjectCloudService/pkg/errcode"
	"WowjoyProject/ObjectCloudService/pkg/general"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

//var token string

// 封装对象相关操作
type Object struct {
	InstanceKey int64
	ResId       string
	Key         string
	Tags        map[string]string
	Path        string
	Count       int
}

func NewObject(data global.ObjectData) *Object {
	return &Object{
		InstanceKey: data.InstanceKey,
		ResId:       global.ObjectSetting.OBJECT_ResId,
		Key:         data.Key,
		Path:        data.Path,
		Count:       data.Count,
	}
}

// 上传对象[POST]
func (obj *Object) UploadObject() {
	global.Logger.Info("开始上传对象：", *obj)
	// 判断文件大小，来区别是否开始分段上传
	// var code string
	// fileSize := general.GetFileSize(obj.Path)
	// if fileSize >= (int64(global.ObjectSetting.File_Fragment_Size << 20)) {
	// 	code = UploadLargeFile(obj, fileSize)
	// } else {
	code := UploadSmallFile(obj)
	if code == "00000" {
		//上传成功更新数据库
		global.Logger.Info("数据上传成功", obj.InstanceKey)
		model.UpdateUplaode(obj.InstanceKey, obj.Key, true)
	} else {
		global.Logger.Info("数据上传失败", obj.InstanceKey)
		// 上传失败时先补偿操作，补偿操作失败后才更新数据库
		if !ReDo(obj, global.UPLOAD) {
			global.Logger.Info("数据补偿失败", obj.InstanceKey)
			// 上传失败更新数据库
			model.UpdateUplaode(obj.InstanceKey, obj.Key, false)
		}
	}
}

// 下载对象[GET]
func (obj *Object) DownObject() {
	// if token == "" {
	// 	// 获取token
	// 	token = "Bearer " + GetToken()
	// }
	// 请求处理太快，http资源没来得及关闭
	// time.Sleep(50 * time.Millisecond)
	global.Logger.Info("开始下载对象：", *obj)
	flag := DownFile(obj)
	if flag {
		global.Logger.Info("下载成功：" + obj.Path)
		// model.UpdateDown(obj.InstanceKey, obj.Key, true)
	} else {
		// 下载失败时先补偿操作，补偿操作失败后才更新数据库
		if !ReDo(obj, global.DOWNLOAD) {
			global.Logger.Info("数据补偿失败", obj.InstanceKey)
			// 下载失败更新数据库
			// model.UpdateDown(obj.InstanceKey, obj.Key, false)
		}
	}
}

// // UploadLargeFile 上传大文件
// func UploadLargeFile(obj *Object, size int64) string {
// 	global.Logger.Debug("开始执行大文件上传", obj.Key)
// 	// num := math.Ceil(float64(size) / float64(global.ObjectSetting.Each_Section_Size))
// 	// 1.初始化
// 	UploadId := Multipart_Upload_Init(obj)
// 	if UploadId == "" {
// 		global.Logger.Error("分段上传初始化获取UploadId是空,结束任务")
// 		return ""
// 	}
// 	global.Logger.Info("UploadId: ", UploadId)
// 	// 2.开始上传小段对象
// 	if Multipart_Upload(obj, UploadId) {
// 		// 文件上传成功完结操作
// 		Multipart_Completion(obj, UploadId)
// 	} else {
// 		// 文件上传失败取消操作
// 		Multipart_Abortion(obj, UploadId)
// 	}
// 	return ""
// }

// UploadSmallFile 上传小文件
func UploadSmallFile(obj *Object) string {
	global.Logger.Debug("开始执行小文件上传")
	// params["resId"] = obj.ResId
	// params["key"] = obj.Key
	url := global.ObjectSetting.OBJECT_POST_Upload
	url += "//"
	url += obj.ResId
	url += "//"
	url += obj.Key
	global.Logger.Debug("操作的URL: ", url)
	file, err := os.Open(obj.Path)
	if err != nil {
		return errcode.File_OpenError.Msg()
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile("file", obj.Path)
	if err != nil {
		global.Logger.Error("CreateFormFile err :", err, file)
		return errcode.Http_HeadError.Msg()
	}
	_, err = io.Copy(formFile, file)
	if err != nil {
		return errcode.File_CopyError.Msg()
	}

	writer.Close()
	// if token == "" {
	// 	// 获取token
	// 	token = "Bearer " + GetToken()
	// }
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		global.Logger.Error("NewRequest err: ", err, url)
		return errcode.Http_RequestError.Msg()
	}
	// request.Header.Set("Authorization", token)
	// 设置AK
	request.Header.Set("accessKey", global.ObjectSetting.OBJECT_AK)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Connection", "close")
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
	}
	resp, err := client.Do(request)
	if err != nil {
		// token = ""
		global.Logger.Error("Do Request got err: ", err)
		return errcode.Http_RequestError.Msg()
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errcode.Http_RespError.Msg()
	}
	global.Logger.Info("resp.Body: ", string(content))
	var result = make(map[string]interface{})
	err = json.Unmarshal(content, &result)
	if err != nil {
		global.Logger.Error("resp.Body: ", "错误")
		return errcode.Http_RespError.Msg()
	}
	// 解析json
	if vCode, ok := result["code"]; ok {
		resultcode := vCode.(string)
		global.Logger.Info("resultcode: ", resultcode)
		return resultcode
	}
	return ""
}

func GetToken() string {
	req, err := http.NewRequest("POST", global.ObjectSetting.TOKEN_URL, nil)
	if err != nil {
		global.Logger.Error("Token NewRequest err:", err, global.ObjectSetting.TOKEN_URL)
		return ""
	}
	req.SetBasicAuth(global.ObjectSetting.TOKEN_Username, global.ObjectSetting.TOKEN_Password)
	req.Header.Set("Connection", "close")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		global.Logger.Error("Do Request got err: ", err, req)
		return ""
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	global.Logger.Info("resp Body: ", string(content))
	var result = make(map[string]interface{})
	_ = json.Unmarshal(content, &result)
	code := result["access_token"]
	var token string
	switch code.(type) {
	case string:
		token = code.(string)
	}
	global.Logger.Info("token: ", token)
	return token
}

// 补偿操作
func ReDo(obj *Object, action global.ActionType) bool {
	global.Logger.Info("开始补偿操作：", obj.InstanceKey)
	if obj.Count < global.ObjectSetting.OBJECT_Count {
		obj.Count += 1
		data := global.ObjectData{
			InstanceKey: obj.InstanceKey,
			Key:         obj.Key,
			Type:        action,
			Path:        obj.Path,
			Count:       obj.Count,
		}
		global.ObjectDataChan <- data
		return true
	}
	return false
}

// DownFile 下载小文件
func DownFile(obj *Object) bool {
	url := global.ObjectSetting.OBJECT_GET_Download
	url += "//"
	url += obj.ResId
	url += "//"
	url += obj.Key
	global.Logger.Debug("操作的URL: ", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		global.Logger.Error("文件下载失败", err, obj.Key)
		return false
	}
	// req.Header.Set("Authorization", token)
	// add params
	// 设置AK
	req.Header.Set("accessKey", global.ObjectSetting.OBJECT_AK)

	resp, err := global.HttpClient.Do(req)
	if err != nil {
		// token = ""
		global.Logger.Error(err)
		return false
	}
	defer resp.Body.Close()
	code := resp.StatusCode
	if code != 200 {
		global.Logger.Error("下载失败：" + obj.Path)
		return false
	}

	len, _ := strconv.ParseInt(resp.Header.Get("Content-size"), 10, 64)
	global.Logger.Info("获取的文件长度：", len)

	general.CheckPath(obj.Path)
	file, _ := os.Create(obj.Path)
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		global.Logger.Error("下载失败：文件拷贝失败：" + obj.Path)
		file.Close()
		os.Remove(obj.Path)
		return false
	} else {
		size := general.GetFileSize(obj.Path)
		global.Logger.Info("下载文件获取的长度：", size)
		if size != len {
			global.Logger.Error("下载失败：保存的文件大小错误：" + obj.Path)
			file.Close()
			os.Remove(obj.Path)
			return false
		}
	}
	return true
}

// // 1.文件分段上传初始化
// func Multipart_Upload_Init(obj *Object) string {
// 	global.Logger.Debug("文件分段上传初始化")
// 	url := global.ObjectSetting.OBJECT_Multipart_Init_URL
// 	url += "//"
// 	url += obj.ResId
// 	url += "//"
// 	url += obj.Key
// 	reqParams := make(map[string]string)
// 	headers := make(map[string]string)
// 	headers["accessKey"] = global.ObjectSetting.OBJECT_AK
// 	headers["Connection"] = "close"
// 	var response = general.PostJson(url, reqParams, headers)
// 	global.Logger.Debug("文件分段上传初始化response", response)
// 	var result = make(map[string]interface{})
// 	err := json.Unmarshal([]byte(response), &result)
// 	if err != nil {
// 		global.Logger.Error("resp.Body: ", "错误")
// 		return ""
// 	}
// 	// 解析json
// 	if vCode, ok := result["code"]; ok {
// 		resultcode := vCode.(string)
// 		if resultcode != "00000" {
// 			global.Logger.Error("文件分段上传初始化接口返回错误", response)
// 			return ""
// 		}
// 	}
// 	if vData, ok := result["data"]; ok {
// 		dataMap := vData.(map[string]interface{})
// 		uploadId := dataMap["uploadId"].(string)
// 		return uploadId
// 	}
// 	return ""
// }

// // 2.分段对象上传
// func Multipart_Upload(obj *Object, uploadid string) bool {
// 	global.Logger.Info("开始执行分段上传函数")
// 	// 将大文件分成小文件
// 	status := true
// 	size := global.ObjectSetting.Each_Section_Size << 20
// 	var fileMap = make(map[int]string)
// 	fileMap = general.FileSplit(obj.Path, size)
// 	global.Logger.Debug("文件分段的map: ", fileMap)
// 	num := len(fileMap)
// 	for v, k := range fileMap {
// 		var code string
// 		var index int
// 		if v == num {
// 			index, code = Multipart_Unifile(obj, k, uploadid, size, num, true)
// 		} else {
// 			index, code = Multipart_Unifile(obj, k, uploadid, size, num, false)
// 		}
// 		if code == "00000" {
// 			//上传成功更新数据库
// 			global.Logger.Info("第", index, "段数据上传成功")
// 			// model.UpdateUplaode(obj.InstanceKey, obj.Key, true)
// 		} else {
// 			global.Logger.Info("第", index, "段数据上传失败")
// 			// model.UpdateUplaode(obj.InstanceKey, obj.Key, false)
// 			status = false
// 		}
// 		os.Remove(k)
// 	}
// 	return status
// }

// // 分段单文件处理
// func Multipart_Unifile(obj *Object, filepath string, uploadid string, size int64, num int, flag bool) (int, string) {
// 	global.Logger.Debug("文件分段上传单文件")
// 	url := global.ObjectSetting.OBJECT_Multipart_Upload_URL
// 	url += "//"
// 	url += obj.ResId
// 	url += "//"
// 	url += obj.Key
// 	file, err := os.Open(filepath)
// 	if err != nil {
// 		return num, errcode.File_OpenError.Msg()
// 	}
// 	defer file.Close()
// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	writer.WriteField("resId", obj.ResId)
// 	writer.WriteField("key", obj.Key)
// 	writer.WriteField("uploadId", uploadid)
// 	writer.WriteField("filePosition", fmt.Sprintf("%d", int64(num-1)*size))
// 	writer.WriteField("partNumber", fmt.Sprintf("%d", num))
// 	if flag {
// 		writer.WriteField("lastPart", "true")
// 	}
// 	formFile, err := writer.CreateFormFile("file", filepath)
// 	if err != nil {
// 		global.Logger.Error("CreateFormFile err :", err, file)
// 		return num, errcode.Http_HeadError.Msg()
// 	}
// 	_, err = io.Copy(formFile, file)
// 	if err != nil {
// 		global.Logger.Error("io.Copy err :", err, file)
// 		return num, errcode.File_CopyError.Msg()
// 	}
// 	writer.Close()
// 	request, err := http.NewRequest("POST", url, body)
// 	if err != nil {
// 		global.Logger.Error("NewRequest err: ", err, url)
// 		return num, errcode.Http_RequestError.Msg()
// 	}
// 	// request.Header.Set("Authorization", token)
// 	// 设置AK
// 	request.Header.Set("accessKey", global.ObjectSetting.OBJECT_AK)
// 	request.Header.Set("Content-Type", writer.FormDataContentType())
// 	request.Header.Set("Connection", "close")
// 	transport := http.Transport{
// 		DisableKeepAlives: true,
// 	}
// 	client := &http.Client{
// 		Transport: &transport,
// 	}
// 	resp, err := client.Do(request)
// 	if err != nil {
// 		global.Logger.Error("Do Request got err: ", err)
// 		return num, errcode.Http_RequestError.Msg()
// 	}
// 	defer resp.Body.Close()
// 	content, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		global.Logger.Error("ioutil.ReadAll got err: ", err)
// 		return num, errcode.Http_RespError.Msg()
// 	}
// 	global.Logger.Info("resp.Body: ", string(content))
// 	var result = make(map[string]interface{})
// 	err = json.Unmarshal(content, &result)
// 	if err != nil {
// 		global.Logger.Error("resp.Body: ", "错误")
// 		return num, errcode.Http_RespError.Msg()
// 	}
// 	// 解析json
// 	if vCode, ok := result["code"]; ok {
// 		resultcode := vCode.(string)
// 		global.Logger.Info("resultcode: ", resultcode)
// 		return num, resultcode
// 	}
// 	return num, errcode.Http_RespError.Msg()
// }

// // 完成对象分段上传
// func Multipart_Completion(obj *Object, uploadid string) string {
// 	global.Logger.Debug("开始执行完成对象分段上传")
// 	url := global.ObjectSetting.OBJECT_Multipart_Completion_URL
// 	url += "//"
// 	url += obj.ResId
// 	url += "//"
// 	url += obj.Key
// 	reqParams := make(map[string]string)
// 	headers := make(map[string]string)
// 	headers["accessKey"] = global.ObjectSetting.OBJECT_AK
// 	headers["Connection"] = "close"
// 	var response = general.PostJson(url, reqParams, headers)
// 	var result = make(map[string]interface{})
// 	err := json.Unmarshal([]byte(response), &result)
// 	if err != nil {
// 		global.Logger.Error("resp.Body: ", "错误")
// 		return ""
// 	}
// 	// 解析json
// 	if vCode, ok := result["code"]; ok {
// 		resultcode := vCode.(string)
// 		global.Logger.Info("resultcode: ", resultcode)
// 		return resultcode
// 	}
// 	return ""
// }

// // 取消对象分段上传
// func Multipart_Abortion(obj *Object, uploadid string) string {
// 	global.Logger.Debug("开始执行完成对象分段上传")
// 	url := global.ObjectSetting.OBJECT_Multipart_Abortion_URL
// 	url += "//"
// 	url += obj.ResId
// 	url += "//"
// 	url += obj.Key
// 	reqParams := make(map[string]string)
// 	headers := make(map[string]string)
// 	headers["accessKey"] = global.ObjectSetting.OBJECT_AK
// 	headers["Connection"] = "close"
// 	var response = general.PostJson(url, reqParams, headers)
// 	var result = make(map[string]interface{})
// 	err := json.Unmarshal([]byte(response), &result)
// 	if err != nil {
// 		global.Logger.Error("resp.Body: ", "错误")
// 		return ""
// 	}
// 	// 解析json
// 	if vCode, ok := result["code"]; ok {
// 		resultcode := vCode.(string)
// 		global.Logger.Info("resultcode: ", resultcode)
// 		return resultcode
// 	}
// 	return ""
// }
