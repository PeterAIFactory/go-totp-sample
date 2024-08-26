package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

type User struct {
	Username string `json:"username"`
	Secret   string `json:"secret"`
}

const dataDir = "./data"

const dbFile = "./users.db"

func main() {
	db, err := sql.Open("sqltes3", dbFile)
	for {
		fmt.Println("\n1. 註冊新用戶")
		fmt.Println("2. 驗證 TOTP")
		fmt.Println("3. 退出")
		fmt.Print("請選擇操作: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			registerUser()
		case 2:
			validateTOTP()
		case 3:
			fmt.Println("再見！")
			return
		default:
			fmt.Println("無效的選擇，請重試。")
		}
	}
}

func registerUser() {
	var username string
	fmt.Print("請輸入用戶名: ")
	fmt.Scanln(&username)

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "YourApp",
		AccountName: username,
	})
	if err != nil {
		log.Fatal(err)
	}

	user := User{
		Username: username,
		Secret:   key.Secret(),
	}

	saveUser(user)

	// 生成 QR 碼
	qrCode, err := qrcode.New(key.URL(), qrcode.Medium)
	if err != nil {
		log.Fatal(err)
	}

	qrFilename := filepath.Join(dataDir, username+"_qr.png")
	err = qrCode.WriteFile(256, qrFilename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("用戶 %s 已註冊。QR 碼已保存為 %s\n", username, qrFilename)
	fmt.Printf("密鑰: %s\n", key.Secret())
}

func validateTOTP() {
	var username string
	fmt.Print("請輸入用戶名: ")
	fmt.Scanln(&username)

	user, err := loadUser(username)
	if err != nil {
		fmt.Println("用戶不存在")
		return
	}

	var code string
	fmt.Print("請輸入 6 位 TOTP 碼: ")
	fmt.Scanln(&code)

	valid := totp.Validate(code, user.Secret)
	if valid {
		fmt.Println("TOTP 驗證成功！")
	} else {
		fmt.Println("TOTP 驗證失敗。")
	}
}

func saveUser(user User) {
	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	filename := filepath.Join(dataDir, user.Username+".json")
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func loadUser(username string) (User, error) {
	filename := filepath.Join(dataDir, username+".json")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
