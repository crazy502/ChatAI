package mail

import (
	"fmt"

	"server/infra/config"

	"gopkg.in/gomail.v2"
)

const (
	CodeMsg     = "GopherAI 验证码如下（2 分钟内有效）："
	UserNameMsg = "GopherAI 为你分配的账号如下，请妥善保存："
)

func SendCaptcha(email, code, msg string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.GetConfig().Email)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "来自 GopherAI 的消息")
	m.SetBody("text/plain", msg+" "+code)

	d := gomail.NewDialer("smtp.qq.com", 587, config.GetConfig().Email, config.GetConfig().Authcode)
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("DialAndSend err %v\n", err)
		return err
	}

	fmt.Println("send mail success")
	return nil
}
