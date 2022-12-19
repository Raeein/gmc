package gmc

import (
	"errors"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Section struct {
	Course Course `json:"course"`
	Code   string `json:"code"`
	Term   string `json:"term"`
}

type Course struct {
	Department string `json:"department"`
	Code       int    `json:"code"`
}

func (c Course) Valid() error {
	if c.Department == "" {
		return errors.New("department cannot be empty")
	}
	if c.Code <= 999 && c.Code >= 9999 {
		return errors.New("course code must be 4 digits")
	}

	return nil
}
