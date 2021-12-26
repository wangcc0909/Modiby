package mail

import (
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"net/smtp"
)

func SendEmail(to, title, content string) {
	host := emailConf.host
	port := emailConf.port
	name := emailConf.name
	from := emailConf.from
	password := emailConf.password
	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s<%s>", name, from)
	header["To"] = to
	header["Subject"] = title
	header["Content-Type"] = "text/html; charset=UTF-8"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + content
	auth := smtp.PlainAuth("", from, password, host)
	err := sendEmailUsingTLS(fmt.Sprintf("%s:%d", host, port),
		auth,
		from,
		[]string{to},
		[]byte(message))
	if err != nil {
		panic(err)
	}
}

func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		logrus.Error("tls.Dial error: ", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func sendEmailUsingTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	//create smtp client
	c, err := dial(addr)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err := c.Auth(auth); err != nil {
				logrus.Error(err)
				return err
			}
		}
	}
	if err := c.Mail(from); err != nil {
		return err
	}

	for _, add := range to {
		if err := c.Rcpt(add); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	return c.Quit()
}
