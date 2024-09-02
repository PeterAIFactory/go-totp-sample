package routes

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

func getTenantTestPage(c *gin.Context) {
	c.HTML(http.StatusOK, "test.html", gin.H{"Key": "Client"})
}

func generateTOTP(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	name, exists := req["name"]
	if !exists {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Missing 'name' parameter",
		})
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "client",
		AccountName: name,
	})
	if err != nil {
		log.Printf("Error generated TOTP: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to generate TOTP",
		})
		return
	}

	var buf bytes.Buffer
	img, err := key.Image(256, 256)
	if err != nil {
		log.Printf("Error generating TOTP image: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to generate TOTP image",
		})
		return
	}

	png.Encode(&buf, img)
	qrCode := base64.StdEncoding.EncodeToString(buf.Bytes())

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "data:image/png;base64," + qrCode,
	})
}
