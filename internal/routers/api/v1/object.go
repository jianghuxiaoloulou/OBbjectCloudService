package v1

import (
	"WowjoyProject/ObjectCloudService/global"
	"WowjoyProject/ObjectCloudService/internal/model"
	"WowjoyProject/ObjectCloudService/pkg/app"
	"WowjoyProject/ObjectCloudService/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// 单文件上传
func UploadFile(c *gin.Context) {
	// 成功：
	app.NewResponse(c).ToResponse("")
	// 失败：
	app.NewResponse(c).ToErrorResponse(errcode.ServerError)
}

// 检查号上传
func UploadNumbers(c *gin.Context) {
	id := c.Param("AccessNumber")
	global.Logger.Info("需要上传的检查号是：", id)
	if id != "" {
		// 获取上传任务：
		model.GetObjectData(id, global.UPLOAD)
		// 成功：
		app.NewResponse(c).ToResponse("")
	} else {
		// 失败：
		app.NewResponse(c).ToErrorResponse(errcode.ServerError)
	}
}

// 单文件下载
func DownFile(c *gin.Context) {
	// 成功：
	app.NewResponse(c).ToResponse("")
	// 失败：
	app.NewResponse(c).ToErrorResponse(errcode.ServerError)
}

// 检查号下载
func DownNumbers(c *gin.Context) {
	id := c.Param("AccessNumber")
	global.Logger.Info("需要下载的检查号是：", id)
	if id != "" {
		// 获取下载任务：
		model.GetObjectData(id, global.DOWNLOAD)
		// 成功：
		app.NewResponse(c).ToResponse("")
	} else {
		// 失败：
		app.NewResponse(c).ToErrorResponse(errcode.ServerError)
	}
}

// 单文件删除
func DeleteFile(c *gin.Context) {
	// 成功：
	app.NewResponse(c).ToResponse("")
	// 失败：
	app.NewResponse(c).ToErrorResponse(errcode.ServerError)
}

// 检查号删除
func DeleteNumbers(c *gin.Context) {
	id := c.Param("AccessNumber")
	global.Logger.Info("需要删除的检查号是：", id)
	if id != "" {
		// 获取下载任务：
		model.GetObjectData(id, global.DELETE)
		// 成功：
		app.NewResponse(c).ToResponse("")
	} else {
		// 失败：
		app.NewResponse(c).ToErrorResponse(errcode.ServerError)
	}
}
