package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/skida12138/drawffxiv-se/services"
)

const sessionName = "DRAWSESSIONXIV"

func registerGlobalMiddlewares(router *gin.Engine) {
	router.Use(setSession())
}

func setSession() func(*gin.Context) {
	return func(context *gin.Context) {
		var userSession services.Session
		var err error
		isSessionValid := true
		defer func() {
			if err != nil {
				throwError(context, err)
			}
		}()
		if temp, err := context.Cookie(sessionName); err == nil {
			userSession = services.Session(temp)
			if isSessionValid, err = userSession.IsValid(); err != nil {
				return
			}
		} else {
			isSessionValid = false
		}
		if !isSessionValid {
			if userSession, err = services.RegisterNewUser(); err != nil {
				return
			}
		}
		context.Set(sessionName, userSession)
		context.SetCookie(sessionName, string(userSession), 86400, "/", "", false, true)
		context.Next()
	}
}

func getSession(context *gin.Context) services.Session {
	session, _ := context.Get(sessionName)
	return session.(services.Session)
}
