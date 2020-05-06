package notify

import (
	"errors"
	"github.com/astaxie/beego"
	"net/smtp"
	"strings"
	"time"
)

type PEmailConfig struct {
	Host string//邮件类型
	Port string//端口号
	User string//用户名
	Pwd string//密码
	From string//谁发送的
}


//定义邮件结构体
type PEmail struct {
	Subject string//邮件主题
	Body string//邮件内容
	To string//发送给谁
	Format string //邮件解析格式

	Config *PEmailConfig//邮件配置文件
}

var (
	config *PEmailConfig
	mainChan chan *PEmail
)


func init() {
	host := beego.AppConfig.String("email.host")
	port := beego.AppConfig.String("email.port")
	from := beego.AppConfig.String("email.from")
	user := beego.AppConfig.String("email.user")
	password := beego.AppConfig.String("email.password")

	poolSize, _ := beego.AppConfig.Int("email.pool")


	config = &PEmailConfig{
		Host:host,
		From:from,
		Port:port,
		User:user,
		Pwd:password,
	}

	//创建通道，用于存储邮件，该通道可以理解为缓冲区
	mainChan = make(chan *PEmail, poolSize)

	go func() {
		for {
			select {
			case m, ok := <- mainChan:
				//判断邮件管道是否关闭
				if !ok {
					return
				}
				//发送邮件
				if err := m.SendToEmail(); err != nil {
					beego.Error("SendEmail:", err.Error())
				}
			}
		}
	}()

}

//发送邮件
func (pe *PEmail) SendToEmail() error {
	//配置用户名，密码等内容
	auth := smtp.PlainAuth("", pe.Config.User, pe.Config.Pwd, pe.Config.Host)

	//切割管理员邮箱
	sendTo := strings.Split(pe.To, ";")

	//获取邮件解析类型
	contentType := GetContentTypeString(pe.Format)
	msg := []byte("To: " + pe.To + "\r\nFrom: " + pe.Config.User +
		"\r\nSubject: " + pe.Subject + "\r\n" + contentType + "\r\n" + pe.Body)

	var err error
	if pe.Config.Port == "25" {
		//发送邮件
		err = smtp.SendMail(pe.Config.Host+":"+pe.Config.Port, auth, pe.Config.User, sendTo, msg)
	}else {
		err = errors.New("邮件发送失败!")
	}
	return err
}


//返回邮件的解析类型
func GetContentTypeString(format string) string {
	var contentType string
	//判断解析类型是否为空
	if format == "" {
		//text/palin:将文件设置为存文本的形式
		contentType = "Content-Type: text/palin" + "; charset=UTF-8"
	}else {
		contentType = "Content-Type: " + format + "; charset=UTF-8"
	}
	return contentType
}


func SendToChan(to, subject, body, mailtype string) bool {
	email := &PEmail{
		Config:config,
		Body: body,
		Subject:subject,
		Format:mailtype,
		To:to,
	}
	select {
	//将邮件发送管道
	case mainChan <- email:
		return true
	//超时控制
	case <- time.After(time.Second * 3):
		return false
	}
}





































