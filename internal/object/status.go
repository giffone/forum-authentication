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

func (s *Statuses) StatusByText(err error, text string, args ...any) {
	sts := statusByText(err, text, args)
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

func StatusByText(err error, text string, args ...any) *Statuses {
	return statusByText(err, text, args)
}

func statusByText(err error, text string, args []any) *Statuses {
	sts := new(Statuses)
	if err != nil {
		e := err.Error()
		log.Printf("status by text: err: %s\n", e)
	}
	if len(args) == 0 {
		sts.StatusBody = text
	} else {
		sts.StatusBody = fmt.Sprintf(text, args...)
	}
	if text == constant.StatusOK {
		sts.StatusCode = constant.Code200
	} else if text == constant.StatusCreated {
		sts.StatusCode = constant.Code201
	} else {
		sts.StatusCode = constant.Code403
	}
	return sts
}
