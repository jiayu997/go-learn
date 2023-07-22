package kapi

import (
	"log"
	"net"
	"net/http"

	"gokit/pkg/webssh"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Webssh(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "webssh.html", nil)
}

func Websshdata(ctx *gin.Context) {
	username := ctx.GetString("username")
	password := ctx.GetString("password")

	id := ctx.Request.Header.Get("Sec-WebSocket-Key")
	addr := ctx.Request.URL.Query().Get("addr")
	wssh := webssh.NewWebSSH()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "ssh connect error")
	}
	wssh.AddSSHConn(id, conn, username, password)
	ws, err := websocket.Upgrade(ctx.Writer, ctx.Request, nil, 1024, 1024)
	if err != nil {
		log.Fatal(err.Error())
	}
	wssh.AddWebsocket(id, ws, username, password)
}
