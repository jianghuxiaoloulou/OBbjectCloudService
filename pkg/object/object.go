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

var token string

// 封装对象相关操作
type Object struct {
	InstanceKey  int64
	BucketId     string
	SyncStrategy string
	Key          string
	Tags         map[string]string
	Path         string
}

func NewObject(data global.ObjectData) *Object {
	var tags = make(map[string]string)
	tags["tag1"] = "test"
	tags["tag2"] = "shulan"
	return &Object{
		InstanceKey:  data.InstanceKey,
		BucketId:     global.ObjectSetting.OBJECT_BucketId,
		SyncStrategy: data.SyncStrategy,
		Key:          data.Key,
		Tags:         tags,
		Path:         data.Path,
	}
}

// 上传对象[POST]
func (obj *Object) UploadObject() {
	global.Logger.Info("开始上传对象：", *obj)
	tag_json, _ := json.Marshal(obj.Tags)
	tag_string := string(tag_json)
	params := make(map[string]string)
	params["bucketId"] = obj.BucketId
	params["syncStrategy"] = obj.SyncStrategy
	params["key"] = obj.Key
	params["tags"] = tag_string
	code := UploadFile(obj.InstanceKey, global.ObjectSetting.OBJECT_POST_Upload, params, "file", obj.Path)
	if code == errcode.Http_Success.Code() {
		//上传成功更新数据库
		global.Logger.Info("数据上传成功", obj.InstanceKey)
		model.UpdateUplaode(obj.InstanceKey, obj.Key, true)
	}
	if code == errcode.Http_Error.Code() {
		// 服务错误不做等待服务重启
		global.Logger.Error("请求错误，等待服务重启")
	}
	if code != errcode.Http_Success.Code() && code != errcode.Http_Error.Code() {
		// 上传失败更新数据库
		global.Logger.Info("数据上传失败", obj.InstanceKey)
		model.UpdateUplaode(obj.InstanceKey, obj.Key, false)
	}
}

// 下载对象[GET]
func (obj *Object) DownObject() {
	if token == "" {
		// 获取token
		token = "Bearer " + GetToken()
	}

	global.Logger.Info("开始下载对象：", *obj)
	params := make(map[string]string)
	params["bucketId"] = obj.BucketId
	params["key"] = obj.Key
	req, err := http.NewRequest(http.MethodGet, global.ObjectSetting.OBJECT_GET_Download, nil)
	if err != nil {
		global.Logger.Error("文件下载失败", err, obj.Key)
		return
	}
	req.Header.Set("Authorization", token)
	// add params
	que := req.URL.Query()
	if params != nil {
		for key, val := range params {
			que.Add(key, val)
		}
		req.URL.RawQuery = que.Encode()
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		token = ""
		global.Logger.Error(err)
		return
	}
	defer resp.Body.Close()
	code := resp.StatusCode
	if code != 200 {
		global.Logger.Error("下载失败：" + obj.Path)
		global.Logger.Error(resp)
		// model.UpdateDown(obj.InstanceKey, obj.Key, false)
		return
	}

	len, _ := strconv.ParseInt(resp.Header.Get("Content-size"), 10, 64)
	global.Logger.Info("获取的文件长度：", len)

	general.CheckPath(obj.Path)
	file, _ := os.Create(obj.Path)
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	size := general.GetFileSize(obj.Path)
	global.Logger.Info("下载文件获取的长度：", size)
	if err != nil {
		global.Logger.Error("文件拷贝失败：" + obj.Path)
		os.Remove(obj.Path)
		return
	} else {
		if size != len {
			global.Logger.Error("保存的文件大小错误：" + obj.Path)
			os.Remove(obj.Path)
			return
		} else {
			global.Logger.Info("下载成功：" + obj.Path)
			// model.UpdateDown(obj.InstanceKey, obj.Key, true)
		}
	}
}

// 删除对象[DELETE]
// func (obj *Object) DelObject() {
// 	req, _ := http.NewRequest("DELETE", setting.OBJECT_DEL_Delete, nil)
// 	res, _ := http.DefaultClient.Do(req)
// 	defer req.Body.Close()
// 	body, _ := ioutil.ReadAll(res.Body)
// 	global.Logger.Debug(string(body))
// }

// 获取对象版本[GET]
// func (obj *Object) GetVersion() {
// 	resp, err := http.Get(setting.OBJECT_GET_Version)
// 	if err != nil {
// 		loggin.Error("获取对象版本错误：", err)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	loggin.Debug(string(body))
// }

// UploadFile 上传文件
func UploadFile(instance_key int64, url string, params map[string]string, paramName, path string) int {
	file, err := os.Open(path)
	if err != nil {
		return errcode.File_OpenError.Code()
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	formFile, err := writer.CreateFormFile(paramName, path)
	if err != nil {
		global.Logger.Error("CreateFormFile err :%v, file: %s", err, file)
		return errcode.Http_HeadError.Code()
	}
	_, err = io.Copy(formFile, file)
	if err != nil {
		return errcode.File_CopyError.Code()
	}

	writer.Close()
	if token == "" {
		// 获取token
		token = "Bearer " + GetToken()
	}
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		global.Logger.Error("NewRequest err: %v, url: %s", err, url)
		return errcode.Http_RequestError.Code()
	}
	request.Header.Set("Authorization", token)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Connection", "close")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		token = ""
		global.Logger.Error("Do Request got err: %v, req: %v", err, request)
		return errcode.Http_RequestError.Code()
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errcode.Http_RespError.Code()
	}
	global.Logger.Info(string(content))
	var result = make(map[string]interface{})
	_ = json.Unmarshal(content, &result)
	code := result["responseCode"]
	var resultcode int
	switch code.(type) {
	case string:
		resultcode, _ = strconv.Atoi(code.(string))
	case int64:
		resultcode = code.(int)
	case float64:
		resultcode = int(code.(float64))
	}
	global.Logger.Info("resultcode: ", resultcode)
	return resultcode
}

func GetToken() string {
	req, err := http.NewRequest("POST", global.ObjectSetting.TOKEN_URL, nil)
	if err != nil {
		global.Logger.Error("Token NewRequest err: %v, url: %s", err, global.ObjectSetting.TOKEN_URL)
		return ""
	}
	req.SetBasicAuth(global.ObjectSetting.TOKEN_Username, global.ObjectSetting.TOKEN_Password)
	req.Header.Set("Connection", "close")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		global.Logger.Error("Do Request got err: %v, req: %v", err, req)
		return ""
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	global.Logger.Info(string(content))
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
