package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"newProject/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	var err error
	DB, err = gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/user?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})

	if err != nil {
		log.Println("gorm conn db failed", err)
	}

	initHandle()

}

// RegisTerUserHandle 注册
func RegisTerUserHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var user model.User
		err = json.Unmarshal(body, &user)
		if err != nil {
			panic(err)
		}

		// 判断用户名或手机号是否存在
		var user1 model.User
		err = DB.Table("user").Where("username = ?", user.Name).First(&user1).Error

		if err != nil {
			log.Println("db find user failed", err)
			w.Write([]byte("failed"))
			return
		}

		if user1.Name == user.Name {
			var m = model.Message{Message: "用户名已存在"}
			jsStr, err := json.Marshal(&m)
			if err != nil {
				log.Println("json marshal failed", err)
				w.Write([]byte("failed"))
				return
			}
			w.Write(jsStr)

		} else if user1.Mobile == user.Mobile {
			var m = model.Message{Message: "手机号已存在"}
			jsStr, err := json.Marshal(&m)
			if err != nil {
				log.Println("json marshal failed", err)
				w.Write([]byte("failed"))
				return
			}
			w.Write(jsStr)
		}

		user.Password = MD5([]byte(user.Password))

		var m model.Message

		err = DB.Table("user").Create(&user).Error
		if err != nil {
			log.Println("db create user failed", err)
			m.Message = "注册失败"
			w.Write([]byte(m.Message))
			return
		}

		m.Message = "注册成功"
		jsStr, err := json.Marshal(&m)
		if err != nil {
			panic(err)
		}
		w.Write(jsStr)

	} else {
		fmt.Fprintf(w, "请求方法错误")
	}
}

// LoginHandle 登录
func LoginHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/json")

	var lr model.LoginRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ReadAll body failed", err)
		w.Write([]byte("failed"))
		return
	}

	err = json.Unmarshal(body, &lr)
	if err != nil {
		log.Println("json Unmarshal body failed", err)
		w.Write([]byte("failed"))
		return
	}

	if lr.Username == "" || lr.Password == "" {

		var m = model.Message{Message: "数据为空"}
		jsStr, err := json.Marshal(&m)
		if err != nil {
			log.Println("json Marshal message struct failed", err)
			w.Write([]byte("failed"))
			return
		}

		w.Write(jsStr)
	} else {

		var user model.User
		err = DB.Table("user").Where("username = ?", lr.Username).First(&user).Error
		if err != nil {
			log.Println("db find user failed", err)
			w.Write([]byte("failed"))
			return
		}

		if MD5([]byte(lr.Password)) == user.Password {
			var m = model.Message{Message: "登录成功"}
			jsStr, err := json.Marshal(&m)
			if err != nil {
				log.Println("json Marshal message struct failed", err)
				w.Write([]byte("failed"))
				return
			}

			w.Write(jsStr)
		} else {
			var m = model.Message{Message: "密码错误"}
			jsStr, err := json.Marshal(&m)
			if err != nil {
				log.Println("json Marshal message struct failed", err)
				w.Write([]byte("failed"))
				return
			}

			w.Write(jsStr)
		}

	}

}

//	initHandle 注册路由
func initHandle() {
	http.HandleFunc("/register", RegisTerUserHandle)
	http.HandleFunc("/login", LoginHandle)
	http.ListenAndServe(":8080", nil)
}

func MD5(str []byte) string {
	m := md5.New()
	m.Write(str)
	c := m.Sum(nil)
	s := hex.EncodeToString(c)
	return s
}

// ChangePasswordHandle 修改用户密码
func ChangePasswordHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/json")

	var cr model.ChangePasswordRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ReadAll body failed", err)
		w.Write([]byte("failed"))
		return
	}

	err = json.Unmarshal(body, &cr)
	if err != nil {
		log.Println("json Unmarshal body failed", err)
		w.Write([]byte("failed"))
		return
	}

	if cr.Username == "" || cr.Password == "" || cr.NewPassword == "" {

		var m = model.Message{Message: "数据为空"}
		jsStr, err := json.Marshal(&m)
		if err != nil {
			log.Println("json Marshal message struct failed", err)
			w.Write([]byte("failed"))
			return
		}

		w.Write(jsStr)
	} else {
		var user model.User
		err = DB.Table("user").Where("username = ?", cr.Username).First(&user).Error
		if err != nil {
			log.Println("db find user failed", err)
			w.Write([]byte("failed"))
			return
		}

		if MD5([]byte(cr.Password)) == user.Password {
			user.Password = MD5([]byte(cr.NewPassword))
			err = DB.Table("user").Save(&user).Error
			if err != nil {
				log.Println("db save user failed", err)
				w.Write([]byte("failed"))
				return
			}

			var m = model.Message{Message: "修改成功"}
			jsStr, err := json.Marshal(&m)
			if err != nil {
				log.Println("json Marshal message struct failed", err)
				w.Write([]byte("failed"))
				return
			}

			w.Write(jsStr)
		} else {
			var m = model.Message{Message: "密码错误"}
			jsStr, err := json.Marshal(&m)
			if err != nil {
				log.Println("json Marshal message struct failed", err)
				w.Write([]byte("failed"))
				return
			}

			w.Write(jsStr)
		}

	}

}
