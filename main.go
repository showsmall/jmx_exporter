package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"sre/jmx_exporter/config"
	"sre/jmx_exporter/controllers"
	"sre/jmx_exporter/gin_logs_func"
	"sre/jmx_exporter/logic"
)

func Loginit() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.SetOutput(os.Stderr)
	log.SetPrefix("jmx_exporter ")
	log.Println("Init logs success")
}

func main() {
	// 初始化日志输出
	Loginit()
	// 初始化配置文件
	err := config.InitConfig()
	if err != nil {
		log.Println("加载配置文件失败:", err)
		return
	}
	log.Println("加载配置文件成功")

	gin.SetMode(config.C.GetString("gin.mode"))
	r := gin.New()
	// 把日志系统集成到gin框架中 把跨域集成到gin
	r.Use(gin_logs_func.GinLogger(), gin_logs_func.GinRecovery(true), gin_logs_func.Cors())
	// 路由
	r.GET("/jmx", logic.Auth, controllers.HandleJmx)

	// 启动服务
	listenportaddr := fmt.Sprintf(":%s", config.C.GetString("host.port"))
	r.Run(listenportaddr)
}
