package tool

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"math/rand"
	"strconv"
	"strings"
	"time"
)


var CodeCache map[string]string

func init() {
	if CodeCache == nil {
		CodeCache = make(map[string]string)
	}
}

//发送邮件
func SendMail(mailTo []string, subject string, body string) error {
	mailConn := map[string]string{
		"user": "zh_account_verify@163.com",
		"pass": "PVQIRHXFMBMFLFEI",//客户端授权密码
		"host": "smtp.163.com",
		"port": "465",
	}
	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int
	m := gomail.NewMessage()
	m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...)    //发送给某个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文
	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	err := d.DialAndSend(m)
	return err
}


//生成有效验证码
func GenValidateCode(width int) string {
	numeric := [10]byte{0,1,2,3,4,5,6,7,8,9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[ rand.Intn(r) ])
	}
	return sb.String()
}