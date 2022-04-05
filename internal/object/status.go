package object

import (
	"fmt"
	"github.com/giffone/forum-authentication/internal/constant"
	"log"
	"net/http"
)

type Statuses struct {
	StatusBody string
	StatusCode int
	ReturnPage string
}

func NewStatuses() *Statuses {
	return new(Statuses)
}

type Status interface {
	Status() *Statuses
}

func (s *Statuses) Status() *Statuses {
	return s
}

func (s *Statuses) StatusByCodeAndLog(code int, err error, message string) {
	log.Printf("statusByCodeAndLog method: %s: %v\n", message, err)
	s.StatusBody = http.StatusText(code)
	s.StatusCode = code
}

func (s *Statuses) StatusByCode(code int) *Statuses {
	s.StatusBody = http.StatusText(code)
	s.StatusCode = code
	return s
}

func (s *Statuses) StatusByText(text, args string, err error) {
	sts := statusByText(text, args, err)
	s.StatusBody = sts.StatusBody
	s.StatusCode = sts.StatusCode
}

func StatusByCodeAndLog(code int, err error, message string) *Statuses {
	log.Printf("%s: %v\n", message, err)
	return &Statuses{
		StatusBody: http.StatusText(code),
		StatusCode: code,
	}
}

func StatusByCode(code int) *Statuses {
	return &Statuses{
		StatusBody: http.StatusText(code),
		StatusCode: code,
	}
}

func StatusByText(text, args string, err error) *Statuses {
	return statusByText(text, args, err)
}

func statusByText(text, args string, err error) *Statuses {
	e := new(Statuses)
	if err != nil {
		log.Printf("err by text: %s\n", err.Error())
	}
	if args == "" {
		e.StatusBody = text
	} else {
		e.StatusBody = fmt.Sprintf(text, args)
	}
	if text == constant.StatusOK {
		e.StatusCode = constant.Code200
	} else if text == constant.StatusCreated {
		e.StatusCode = constant.Code201
	} else {
		e.StatusCode = constant.Code403
	}
	return e
}
