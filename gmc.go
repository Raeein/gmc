package gmc

import "fmt"

type Course struct {
	Year    string `json:"year"`
	Section string `json:"section"`
	ID      string `json:"id"`
}

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (c Course) String() string {
	return fmt.Sprintf("%s, %s, %s\n", c.Year, c.Section, c.ID)
}
