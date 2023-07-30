package kapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 测试接口
func Test(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "test.html", nil)
	fmt.Println(ctx.Request.Host)
}

func Testdata(ctx *gin.Context) {
	//	token := ctx.PostForm("token")
	//	if token == "123" {
	//		ctx.Data(http.StatusOK, klog.ContentTypeHTML, []byte(klog.GenerateHtml(1, "i'am ok", false)))
	//	} else {
	//		ctx.String(http.StatusInternalServerError, "falied")
	//	}
}

func Testlog(ctx *gin.Context) {
	//	upGrader := websocket.Upgrader{
	//		ReadBufferSize:  128,
	//		WriteBufferSize: 128,
	//		CheckOrigin: func(r *http.Request) bool {
	//			return true
	//		},
	//	}
	//	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	//	if err != nil {
	//		ctx.Abort()
	//	}
	//	ty, by, _ := ws.ReadMessage()
	//	fmt.Println(ty, string(by))
	//	for {
	//		time.Sleep(time.Second * 1)
	//	}
	//	defer ws.Close()

}
