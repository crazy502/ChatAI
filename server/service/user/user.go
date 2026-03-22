package user

import (
	"log"

	"server/common/code"
	myemail "server/common/email"
	myredis "server/common/redis"
	"server/dao/user"
	"server/model"
	"server/utils"
	"server/utils/myjwt"
)

func Login(username, password string) (string, bool, code.Code) {
	var userInformation *model.User
	var ok bool

	if ok, userInformation = user.IsExistUser(username); !ok {
		return "", false, code.CodeUserNotExist
	}

	if utils.IsBcryptHash(userInformation.Password) {
		if !utils.VerifyPassword(userInformation.Password, password) {
			return "", false, code.CodeInvalidPassword
		}
	} else {
		if userInformation.Password != utils.MD5(password) {
			return "", false, code.CodeInvalidPassword
		}

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			log.Println("Login hash legacy password error:", err)
		} else if err := user.UpdatePassword(userInformation.ID, hashedPassword); err != nil {
			log.Println("Login migrate legacy password error:", err)
		}
	}

	token, err := myjwt.GenerateToken(userInformation.ID, userInformation.Username, userInformation.IsAdmin)
	if err != nil {
		return "", false, code.CodeServerBusy
	}

	return token, userInformation.IsAdmin, code.CodeSuccess
}

func Register(email, password, captcha string) (string, bool, code.Code) {
	var ok bool
	var userInformation *model.User

	if ok, _ = user.IsExistEmail(email); ok {
		return "", false, code.CodeEmailExist
	}

	if ok, _ = myredis.CheckCaptchaForEmail(email, captcha); !ok {
		return "", false, code.CodeInvalidCaptcha
	}

	username := utils.GetRandomNumbers(11)
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", false, code.CodeServerBusy
	}

	if userInformation, ok = user.Register(username, email, hashedPassword, false); !ok {
		return "", false, code.CodeServerBusy
	}

	if err := myemail.SendCaptcha(email, username, user.UserNameMsg); err != nil {
		return "", false, code.CodeServerBusy
	}

	token, err := myjwt.GenerateToken(userInformation.ID, userInformation.Username, userInformation.IsAdmin)
	if err != nil {
		return "", false, code.CodeServerBusy
	}

	return token, userInformation.IsAdmin, code.CodeSuccess
}

func SendCaptcha(email string) code.Code {
	sendCode := utils.GetRandomNumbers(6)
	if err := myredis.SetCaptchaForEmail(email, sendCode); err != nil {
		return code.CodeServerBusy
	}

	if err := myemail.SendCaptcha(email, sendCode, myemail.CodeMsg); err != nil {
		return code.CodeServerBusy
	}

	return code.CodeSuccess
}
