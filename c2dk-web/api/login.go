package kapi

import (
	"gokit/pkg/jwt"
	"gokit/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 登陆时，校验数据
// method: POST
// interface: /api/v1/login
func Logindata(ctx *gin.Context) {
	username := ctx.PostForm("username")
	token := ctx.PostForm("token")
	if username != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "认证失败"})
		return
	}
	err := jwt.ParseJwtToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "token认证失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "认证成功"})

	// 设置认证通过
	utils.Auth = true
}

// 认证中间件，如果未认证，则后续不会执行
func Auth(ctx *gin.Context) {
	//fmt.Println("当前接口为：", ctx.Request.URL.Path, "登陆状态为：", utils.Auth)
	// login 校验登陆数据接口，放行
	//	if ctx.Request.URL.Path == "/api/v1/login" && ctx.Request.Method == "POST" {
	//		return
	//	}
	//
	//	if utils.Auth {
	//		ctx.Next()
	//	} else {
	//		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
	//		ctx.Abort()
	//	}
}

func Login(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", nil)
}
