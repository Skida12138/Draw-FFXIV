package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/skida12138/drawffxiv-se/services"
	"github.com/skida12138/drawffxiv-se/wsproto"
)

func registerUsersRoutes(router *gin.Engine) {
	router.GET("/users/current", handleGetNickName)
}

func handleGetNickName(context *gin.Context) {
	nickName, err := getSession(context).GetNickName()
	if err != nil {
		nickName = ""
	}
	accepted(context, &gin.H{
		"nickName": nickName,
	})
}

var buff map[services.Session][]byte

func handleConn(conn *wsproto.Conn, session services.Session) {
	for {
		msgType, msg, err := conn.SyncRead()
		if err != nil {
			break
		}
		if msgType == websocket.TextMessage {
			buff[session] = append(buff[session], msg...)
			buff[session], err = services.Consume(buff[session], session)
		}
	}
}
