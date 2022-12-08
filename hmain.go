package main

// go-bindata -o asset/asset.go -pkg asset  website/...

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/MX-Steve/my-go-app/api"
	"github.com/MX-Steve/my-go-app/asset"
	"github.com/MX-Steve/my-go-app/handler"
	"github.com/MX-Steve/my-go-app/middlerware"
	"github.com/MX-Steve/my-go-app/vip"

	"github.com/gin-gonic/gin"
)

func runRouterMain() {
	router := gin.Default()

	router.StaticFS("/public", http.Dir("./website/static"))
	router.StaticFS("/kindeditor", http.Dir("./website/static/kindeditor"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	router.LoadHTMLGlob("website/tpl/*/*")
	v := router.Group("/")
	{
		v.GET("/", handler.IndexHandler)
		v.GET("/login", handler.LoginHandler)
	}
	noAuthGroup := router.Group("/out")
	{
		noAuthGroup.POST("/api/login", api.LoginUser)
	}
	authGroup := router.Group("/in")
	authGroup.Use(middlerware.JWTAuth())
	{
		authGroup.GET("/api/tools/v1/alarm", api.GetAlarmClock)
		authGroup.POST("/api/tools/v1/alarm", api.CudAlarmClock)
		authGroup.GET("/api/tools/v1/audit", api.GetAudits)
		authGroup.GET("/api/tools/v1/auditbyid", api.GetAuditById)
		authGroup.GET("/api/tools/v1/auditbyname", api.GetAuditByName)
		authGroup.POST("/api/tools/v1/plan", api.CudPlans)
		authGroup.GET("/api/tools/v1/plan", api.GetPlans)
		authGroup.GET("/api/tools/v1/planbyid", api.GetPlanById)
		authGroup.POST("/api/tools/v1/article", api.CudArticles)
		authGroup.GET("/api/tools/v1/article", api.GetArticles)
		authGroup.GET("/api/tools/v1/articlebyid", api.GetArticleById)
		authGroup.GET("/api/tools/v1/articlebyname", api.GetArticleByName)
		authGroup.GET("/api/tools/v1/task", api.GetTasks)
		authGroup.POST("/api/tools/v1/task", api.CudTasks)
		authGroup.GET("/api/tools/v1/taskbyid", api.GetTaskById)
		authGroup.GET("/api/userinfo", api.GetUser)
		authGroup.POST("/api/usermod", api.ModUser)
		authGroup.POST("/api/usermodphoto", api.ModPhoto)
		authGroup.POST("/api/sendemail", api.SendEmail)
	}
	port := vip.GetIniData("version.port")
	router.Run(fmt.Sprintf(":%s", port))
}

func runRouterFile() {
	router := gin.Default()
	imgP := vip.GetIniData("router80.imgPath")
	router.StaticFS("/", http.Dir(imgP))

	router.Run(":80")
}
func init_http_server() {
	go runRouterFile()
	runRouterMain()
}

func getWebSite() {
	// 释放 go-bindata 压缩进二进制包的数据
	names := asset.AssetNames()
	for _, name := range names {
		dir, _ := path.Split(name)
		_, exist := os.Stat(dir)
		if os.IsNotExist(exist) {
			os.MkdirAll(dir, os.ModePerm)
		}
		bytes, _ := asset.Asset(name)
		ioutil.WriteFile(name, bytes, 0644)
	}
}
