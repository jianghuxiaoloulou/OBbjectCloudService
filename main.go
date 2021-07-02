package main

import (
	"WowjoyProject/ObjectCloudService/global"
	"WowjoyProject/ObjectCloudService/internal/routers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	readSetup()
}

// @title 对象存储系统
// @version 1.0
// @description 对象文件上传下载
// @termsOfService https://github.com/jianghuxiaoloulou/OBbjectCloudService.git
func main() {
	global.Logger.Info("hello")
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()

	ser := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	ser.ListenAndServe()
}
