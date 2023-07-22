package utils

import (
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func SendEmail(fileName string) {
	em := email.NewEmail()
	em.From = "qujiayu98@163.com"
	em.To = []string{"2281823407@qq.com"}
	em.Subject = fileName
	err := em.Send("smtp.qq.com:22", smtp.PlainAuth("", "自己的邮箱账号", "自己邮箱的授权码", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err.Error())
	}
}
