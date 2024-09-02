package routes

import (
	"Go_learning/model"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

// type User struct {
// 	Username string `json:"username"`
// 	Secret   string `json:"secret"`
// }

const dataDir = "data"

func getTenantTestPage(c *gin.Context) {
	c.HTML(http.StatusOK, "test.html", gin.H{"Key": "Client"})
}

func registerTOTP(c *gin.Context) {
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

	user := model.WebUser{
		Username: name,
		Secret:   key.Secret(),
	}

	// save user's data in jason
	saveUser(user)
	// Generate QRcode
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

func validateTOTP(c *gin.Context) {
	var user model.WebUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	// require database's user
	userInData, err := loadUser(user.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	valid := totp.Validate(user.Secret, userInData.Secret)
	if !valid {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "TOTP validate false",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "TOTP validate successful",
	})
}

func saveUser(user model.WebUser) {
	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	filename := filepath.Join(dataDir, user.Username+".json")
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func loadUser(username string) (model.WebUser, error) {
	filename := filepath.Join(dataDir, username+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return model.WebUser{}, err
	}

	var user model.WebUser
	err = json.Unmarshal(data, &user)
	if err != nil {
		return model.WebUser{}, err
	}

	return user, nil
}
