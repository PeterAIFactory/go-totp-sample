package routes

import (
	"crypto/sha256"
	"encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
)

func InitRoutes(p *gin.Engine) {
	secret := sha256.Sum256([]byte("1qaz2wsx"))
	store := cookie.NewStore(secret[:])
	p.Use(sessions.Sessions("session", store))
	gob.Register(webauthn.SessionData{})

	p.GET("/login", getTenantTestPage)

	p.POST("/api/totp/registration", registerTOTP)

	p.POST("/api/totp/validation", validateTOTP)

	// load templates
	p.LoadHTMLGlob("static/templates/*.html")
}
