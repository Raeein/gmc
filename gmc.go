package gmc

import (
	"context"
	"errors"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (w User) Valid() error {
	if w.Email == "" && w.Name == "" {
		return errors.New("no contact information provided")
	}

	return nil
}

type Section struct {
	Course Course `json:"course"`
	Code   string `json:"code"`
	Term   string `json:"term"`
}

func (s Section) Valid() error {
	if s.Code == "" {
		return errors.New("section code cannot be empty")
	}
	if s.Term == "" {
		return errors.New("term cannot be empty")
	}
	return s.Course.Valid()
}

type Course struct {
	Department string `json:"department"`
	Code       int    `json:"code"`
}

func (c Course) Valid() error {
	if c.Department == "" {
		return errors.New("Department cannot be empty")
	}
	if c.Code <= 999 && c.Code >= 9999 {
		return errors.New("Course code must be 4 digits")
	}

	return nil
}

type Notifier interface {
	Notify(context.Context, Section, ...User) error
}
