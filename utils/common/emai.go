package common

import (
	"gopkg.in/gomail.v2"
)

const (
	USER_NAME        = "" //邮箱
	AUTH_CODE        = "" //授权码
	SEND_SERVER_HOST = "" //stmp发送邮件服务器
	SEND_SERVER_PORT = 0  //stmp发送邮件服务器端口
	SENDER_NAME      = "" //发件人昵称
)

// addressee:收件人
// title:主题
// content:内容
func SendMail(addressee, title, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(USER_NAME, SENDER_NAME))
	m.SetHeader("To", addressee)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)
	d := gomail.NewDialer(SEND_SERVER_HOST, SEND_SERVER_PORT, USER_NAME, AUTH_CODE)
	err := d.DialAndSend(m)
	return err
}
