package email

import (
	"bytes"
	"fmt"
	"github.com/Raeein/gmc"
	"net/smtp"
	"strconv"
	"text/template"
)

func Send(u gmc.User, courses []string, host string, port int, from string, passwd string) {

	auth := smtp.PlainAuth("", from, passwd, host)
	uri := host + ":" + strconv.FormatInt(int64(port), 10)
	text, _ := makeTemplate(u, courses)
	err := smtp.SendMail(uri, auth, from, []string{u.Email}, text.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}

func makeTemplate(u gmc.User, courses []string) (bytes.Buffer, error) {
	var body bytes.Buffer
	var tpl *template.Template

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Spots Availabe \n%s\n\n", mimeHeaders)))

	tpl = template.Must(template.ParseGlob("template/*.gohtml"))

	_ = tpl.Execute(&body, struct {
		Name    string
		Courses []string
	}{
		Name:    u.Name,
		Courses: courses,
	})
	return body, nil
}
