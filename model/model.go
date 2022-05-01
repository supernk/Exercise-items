package model

type Smoke struct {
	Money    int `json:"money"`
	Quantity int `json:"quantity"`
}

type Message struct {
	Message string `json:"message"`
}

type UserRequest struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Mobile string `json:"mobile"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Name     string `json:"name" gorm:"column:username"`
	Mobile   string `json:"mobile" gorm:"column:phone"`
	Password string `json:"password" gorm:"column:password"`
}

type ChangePasswordRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"newpassword"`
}
