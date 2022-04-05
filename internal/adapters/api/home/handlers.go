package home

import (
	"context"
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/model"
	"github.com/giffone/forum-authentication/internal/service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type hHome struct {
	sPost     service.Post
	sCategory service.Category
}

func NewHandler(sPost service.Post, sCategory service.Category) api.Handler {
	return &hHome{
		sPost:     sPost,
		sCategory: sCategory,
	}
}

func (hh *hHome) Register(ctx context.Context, router *http.ServeMux, session api.Session) {
	router.HandleFunc(constant.URLHome, session.Check(ctx, hh.Home))
	router.HandleFunc(constant.URLFavicon, hh.Favicon)
	router.HandleFunc(constant.URLCategoryBy, session.Check(ctx, hh.ByCategory))
}

func (hh *hHome) Home(ctx context.Context, ck *object.Cookie, sts object.Status,
	w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " ", r.URL.Path)
	// check errors in cookie
	if sts != nil {
		api.Message(w, sts)
		return
	}
	if r.Method != "GET" {
		api.Message(w, object.StatusByCode(constant.Code405))
		return
	}
	if r.URL.Path != "/" {
		api.Message(w, object.StatusByCode(constant.Code404))
		return
	}
	posts := model.NewPosts(nil, ck)
	posts.St.AllPost = true
	hh.get(ctx, posts, w)
}

func (hh *hHome) ByCategory(ctx context.Context, ck *object.Cookie, sts object.Status,
	w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " ", r.URL.Path)
	// check errors in cookie
	if sts != nil {
		api.Message(w, sts)
		return
	}
	if r.Method != "GET" {
		api.Message(w, object.StatusByCode(constant.Code405))
		return
	}
	// get id category from url
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, constant.URLCategoryBy))
	if err != nil || id == 0 {
		sts := object.StatusByCodeAndLog(constant.Code400,
			err, "handler: postFormValue: atoi")
		api.Message(w, sts)
		return
	}
	posts := model.NewPosts(nil, ck)
	posts.MakeKeys(constant.KeyCategory, id)
	posts.St.AllPost = true
	hh.get(ctx, posts, w)
}

func (hh *hHome) Favicon(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " ", r.URL.Path)
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	http.ServeFile(w, r, "assets/ico/favicon.ico")
}

func (hh *hHome) get(ctx context.Context, m model.Models, w http.ResponseWriter) {
	ctx, cancel := context.WithTimeout(ctx, constant.TimeLimit)
	defer cancel()
	// parse
	pe, sts := api.NewParseExecute("index").Parse()
	if sts != nil {
		api.Message(w, sts)
		return
	}
	p := m.Return().Posts
	// check session
	pe.Data["Session"] = p.Ck.Session
	// get data
	pe.Data["AllPost"] = p.St.AllPost
	pe.Data["Posts"], sts = hh.sPost.Get(ctx, m)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	c := model.NewCategories(nil, nil)
	pe.Data["Category"], sts = hh.sCategory.GetList(ctx, c)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// execute
	pe.Execute(w, constant.Code200)
}
