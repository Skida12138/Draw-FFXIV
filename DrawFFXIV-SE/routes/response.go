package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skida12138/drawffxiv-se/i18n"
)

func throwError(context *gin.Context, err error) {
	log.Panic(err)
	context.JSON(http.StatusInternalServerError, gin.H{
		"errMsg": i18n.Msg("serverErrorHint"),
	})
}

func badRequest(context *gin.Context, errMsg string) {
	context.JSON(http.StatusBadRequest, gin.H{
		"errMsg": errMsg,
	})
}

func accepted(context *gin.Context, result *gin.H) {
	context.JSON(http.StatusAccepted, result)
}
