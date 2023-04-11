package main

import (
	"WowjoyProject/ObjectCloudService/global"
	"WowjoyProject/ObjectCloudService/internal/model"
	"WowjoyProject/ObjectCloudService/internal/routers"
	"WowjoyProject/ObjectCloudService/pkg/object"
	"WowjoyProject/ObjectCloudService/pkg/workpattern"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @title 对象存储系统
// @version 1.0.0.1
// @description 对象文件上传下载
// @termsOfService https://github.com/jianghuxiaoloulou/OBbjectCloudService.git
func main() {
	global.Logger.Info("*******开始运行对象存储系统********")

	global.ObjectDataChan = make(chan global.ObjectData)

	// 注册工作池，传入任务
	// 参数1 初始化worker(工人)设置最大线程数
	wokerPool := workpattern.NewWorkerPool(global.GeneralSetting.MaxThreads)
	// 有任务就去做，没有就阻塞，任务做不过来也阻塞
	wokerPool.Run()
	// 处理任务
	go func() {
		for {
			select {
			case data := <-global.ObjectDataChan:
				sc := &Dosomething{key: data}
				wokerPool.JobQueue <- sc
			}
		}
	}()

	// 根据工作模式判断
	switch global.GeneralSetting.WorkMode {
	case int(global.AutoMode):
		global.Logger.Info("工作模式配置为自动上传模式")
		auto()
	case int(global.WebMode):
		global.Logger.Info("工作模式配置为Web服务模式")
		web()
	default:
		global.Logger.Info("工作模式配置错误，程序退出")
	}
}

type Dosomething struct {
	key global.ObjectData
}

func (d *Dosomething) Do() {
	global.Logger.Info("正在处理的数据是：", d.key)
	//处理封装对象
	obj := object.NewObject(d.key)
	switch d.key.Type {
	case global.UPLOAD:
		// 数据上传
		obj.UploadObject()
	case global.DOWNLOAD:
		// 数据下载
		obj.DownObject()
	case global.DELETE:
		// 数据删除
		//obj.DelObject()
	}
}

func auto() {
	// 方式一：
	for {
		time.Sleep(time.Second * 10)
		model.AutoUploadObjectData()
	}

	// 方式二：获取任务(定时任务)
	// MyCron := cron.New()
	// MyCron.AddFunc(global.GeneralSetting.CronSpec, func() {
	// 	global.Logger.Info("开始执行定时任务")
	// 	model.AutoUploadObjectData()
	// })

	// MyCron.Start()

	// defer MyCron.Stop()

	// select {}
}

func web() {
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
