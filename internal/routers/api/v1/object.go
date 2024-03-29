package v1

import (
	"WowjoyProject/ObjectCloudService/global"
	"WowjoyProject/ObjectCloudService/internal/model"
	"WowjoyProject/ObjectCloudService/pkg/app"
	"WowjoyProject/ObjectCloudService/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// 检查号上传
func UploadNumbers(c *gin.Context) {
	id := c.Param("AccessNumber")
	global.Logger.Info("需要上传的检查号是：", id)
	if id != "" {
		// 成功：
		app.NewResponse(c).ToResponse(nil)
		// 获取上传任务：
		model.GetObjectData(id, global.UPLOAD)
	} else {
		// 失败：
		app.NewResponse(c).ToErrorResponse(errcode.ServerError)
	}
}

// 检查号下载
func DownNumbers(c *gin.Context) {
	id := c.Param("AccessNumber")
	global.Logger.Info("需要下载的检查号是：", id)
	if id != "" {
		// 成功：
		app.NewResponse(c).ToResponse(nil)
		// 获取下载任务：
		model.GetObjectData(id, global.DOWNLOAD)
	} else {
		// 失败：
		app.NewResponse(c).ToErrorResponse(errcode.ServerError)
	}
}
