package routes

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skida12138/drawffxiv-se/i18n"
	"github.com/skida12138/drawffxiv-se/services"
	"github.com/skida12138/drawffxiv-se/utils"
)

func registerRoomsRoutes(router *gin.Engine) {
	roomRouter := router.Group("/rooms")
	roomRouter.POST("/", handleRoomEntering)
}

func handleRoomEntering(context *gin.Context) {
	var params struct {
		RoomID   string `json:"roomID"`
		NickName string `json:"nickName"`
	}
	userSession := getSession(context)
	throwError(context, utils.NewResult(
		nil, context.BindJSON(&params),
	).AndThen(func(_ interface{}) (interface{}, error) {
		return nil, userSession.SetNickName(params.RoomID)
	}).Error())
	if roomID, err := strconv.ParseUint(params.RoomID, 10, 64); err != nil || roomID > ((1<<31)-1) {
		badRequest(context, i18n.Msg("roomIDShouldBePosInt"))
	} else {
		services.UserEnterRoom(getSession(context), int32(roomID))
	}
}
