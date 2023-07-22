package router

import (
	"fmt"
	kapi "gokit/api"
	"gokit/pkg/jwt"
	"log"
	"net/http"

	"gokit/utils"

	"github.com/gin-gonic/gin"
)

func InitGin() {
	// web 启动前检查
	err := utils.Amp.InitFlag()
	if err != nil {
		log.Fatal(err.Error())
	}

	// 生成GWT配置
	utils.JwtToken, err = jwt.GenerateJwtToken()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("监听端口为：*:8888")
	fmt.Printf("\033[1;37;41m登陆token为：admin/%s\033[0m\n", utils.JwtToken)

	// 初始化gin
	router := gin.Default()

	// 初始化静态资源
	router.LoadHTMLGlob("./static/html/*")

	// 初始化css静态资源
	router.Static("/static", "./static/")

	// 初始路由
	SetUpRoutes(router)

	// 启动gin
	router.Run(":8888")
}

func SetUpRoutes(router *gin.Engine) {
	// web api接口组
	api := router.Group("/api/v1").Use(kapi.Auth)
	{
		// 登陆数据处理
		api.POST("login", kapi.Logindata)

		// 部署前检查提交上来的数据
		api.GET("deploycheck", kapi.DeployCheck)

		// 生成ansible配置文件，并执行Ansible
		api.GET("deployrun", kapi.DeployRun)

		// 部署日志下载接口
		api.GET("log", kapi.DeployLog)

		// 集群异常pod清理接口
		api.DELETE("deletepod", kapi.DeleteErrorPods)

		// 集群DNS延迟检查接口
		api.GET("dnsdelay", kapi.DnsDelay)

		// 集群POD状态接口
		api.GET("podstatus", kapi.PodStatus)

		// 集群删除node节点接口
		api.DELETE("deletenode", kapi.DeleteNode)

		// 集群添加node 节点接口
		api.GET("addnode", kapi.AddNode)

		// 集群SVC 暴露端口查询
		api.GET("service", kapi.GetSVC)

		// 输出平台访问地址
		api.GET("urlacess", kapi.UrlAcess)

		// 单独部署组件
		api.GET("component", kapi.DeployComponent)

	}

	// 登陆接口
	router.GET("/login", kapi.Login)

	// 跳转套登陆页面
	router.GET("/", func(ctx *gin.Context) {
		if ctx.GetBool("auth") {
			ctx.Next()
		} else {
			ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		}
	})

	// 认证成功后，跳转到部署页面
	router.GET("/deploy", kapi.Auth, kapi.Deploy)

	// 测试接口，用于前端测试
	router.GET("/test", kapi.Test)
	router.POST("/test", kapi.Testdata)
	router.GET("/testlog", kapi.Testlog)

	// web ssh接口
	router.GET("/webssh", kapi.Webssh)
	api.GET("/webssh/data", kapi.Websshdata)
}
