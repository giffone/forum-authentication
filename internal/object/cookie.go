package object

import (
	"errors"
	"fmt"
	"github.com/giffone/forum-authentication/internal/constant"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Cookie struct {
	User        int
	Post        int
	PostString  string
	Session     bool
	SessionUUID string
}

func NewCookie() *Cookie {
	return new(Cookie)
}

func (c *Cookie) CookieUserIDRead(r *http.Request) Status {
	ck, err := r.Cookie(constant.CookieUserID)
	if err != nil {
		return StatusByCode(constant.Code400)
	}
	user, err := strconv.Atoi(ck.Value)
	if err != nil || user == 0 {
		return StatusByCodeAndLog(constant.Code400,
			err, "cookie: cookieUserIDRead: atoi")
	}
	c.User = user
	log.Printf("object cookie: userID is: %d", c.User)
	return nil
}

func (c *Cookie) CookieSessionRead(r *http.Request) Status {
	ck, err := r.Cookie(constant.CookieSession)
	if err != nil {
		return StatusByCodeAndLog(constant.Code401,
			err, "cookie: cookieSessionRead")
	}
	c.SessionUUID = ck.Value
	log.Printf("object cookie: middleware-uuid is: %s", c.SessionUUID)
	return nil
}

func (c *Cookie) CookiePostIDRead(r *http.Request) Status {
	p, err := r.Cookie(constant.CookiePostID)
	if err != nil {
		return StatusByCodeAndLog(constant.Code400,
			err, "cookie: cookiePostIDRead")
	}
	c.PostString = p.Value
	log.Printf("object cookie: postID is: %s", c.PostString)
	return nil
}

func (c *Cookie) AddUser(id int) *Cookie {
	c.User = id
	return c
}

func CookieSessionAndUserID(w http.ResponseWriter, value []string, method string) Status {
	name := []string{constant.CookieSession, constant.CookieUserID}
	sts := cookieSet(w, name, value, method)
	if sts != nil {
		return sts
	}
	return nil
}

func CookiePostID(w http.ResponseWriter, id string) Status {
	name := []string{constant.CookiePostID}
	value := []string{id}
	sts := cookieSet(w, name, value, "")
	if sts != nil {
		return sts
	}
	return nil
}

func CookiePostIDDel(w http.ResponseWriter) Status {
	name := []string{constant.CookiePostID}
	value := []string{""}
	sts := cookieSet(w, name, value, "erase")
	if sts != nil {
		return sts
	}
	return nil
}

func cookieSet(w http.ResponseWriter, name []string, value []string, method string) Status {
	lName := len(name)
	lValue := len(value)
	if lName != lValue {
		e := fmt.Sprintf(" different length name(%d) and value(%d)\n", lName, lValue)
		err := errors.New(e)
		return StatusByCodeAndLog(constant.Code500,
			err, "create cookie:")
	}
	for i := 0; i < len(name); i++ {
		c := &http.Cookie{}
		c.Name = name[i]
		c.Value = value[i]
		c.Path = "/"
		if method == "remember" {
			c.Expires = time.Now().AddDate(0, 0, constant.SessionExpire)
			c.MaxAge = constant.SessionMaxAge
		} else if method == "erase" {
			c.Expires = time.Unix(1, 0)
			c.MaxAge = -1
		}
		http.SetCookie(w, c)
	}
	return nil
}
