package email

import (
	"fmt"
	"github.com/Raeein/gmc"
	"net/smtp"
	"strconv"
)

func Send(u gmc.User, host string, port int, from string, passwd string) {

	auth := smtp.PlainAuth("", from, passwd, host)
	uri := host + ":" + strconv.FormatInt(int64(port), 10)

	err := smtp.SendMail(uri, auth, from, []string{u.Email}, []byte(`HIIII`))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}
