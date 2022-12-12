package gmc

import "fmt"

type Sender struct {
	Email    string
	Password string
}

type User struct {
	Email string
	Name  string
}

func (u *User) String() string {
	return fmt.Sprintf("%s: %s\n", u.Email, u.Name)
}

type EmailServer struct {
	SmtpHost string
	SmtpPort string
}
