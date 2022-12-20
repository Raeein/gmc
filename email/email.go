package email

import (
	"bytes"
	"fmt"
	"github.com/Raeein/gmc"
	"net/smtp"
	"text/template"
)

type Service struct {
	host     string
	port     string
	password string
	from     string
}

func New(host string, port, password, from string) Service {
	return Service{
		host:     host,
		port:     port,
		password: password,
		from:     from,
	}
}

func (s Service) Send(section gmc.Section, users ...gmc.User) error {
	auth := smtp.PlainAuth("", s.from, s.password, s.host)
	uri := s.host + ":" + s.port
	var err error
	for _, u := range users {
		content, _ := makeContent(u, section)
		e := smtp.SendMail(uri, auth, s.from, []string{u.Email}, content.Bytes())
		if e != nil {
			fmt.Println(err)
			err = e
		}
	}
	return err
}

func makeContent(u gmc.User, section gmc.Section) (bytes.Buffer, error) {
	var body bytes.Buffer
	var tpl *template.Template
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Spots Availabe \n%s\n\n", mimeHeaders)))

	tpl = template.Must(template.ParseFiles("email.gohtml"))

	_ = tpl.Execute(&body, struct {
		Name    string
		Section gmc.Section
	}{
		Name:    u.Name,
		Section: section,
	})
	return body, nil
}
