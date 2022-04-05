package post

import (
	"context"
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/object/model"
	"github.com/giffone/forum-authentication/internal/service"
	"log"
	"net/http"
	"strconv"
)

type hPost struct {
	service   service.Post
	sCategory service.Category
	sComment  service.Comment
	ratio     api.Ratio
}

func NewHandler(service service.Post,
	sCategory service.Category, sComment service.Comment, ratio api.Ratio) api.Handler {
	return &hPost{
		service:   service,
		sCategory: sCategory,
		sComment:  sComment,
		ratio:     ratio,
	}
}

func (hp *hPost) Register(ctx context.Context, router *http.ServeMux, session api.Session) {
	router.HandleFunc(constant.URLRead, session.Check(ctx, hp.Read))
	router.HandleFunc(constant.URLPost, session.Check(ctx, hp.CreatePost))
	router.HandleFunc(constant.URLComment, session.Check(ctx, hp.CreateComment))
	router.HandleFunc(constant.URLReadRatio, session.Check(ctx, hp.CreateLike))
}

func (hp *hPost) Read(ctx context.Context, ck *object.Cookie, sts object.Status,
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
	u := r.URL.Query()
	ck.PostString = u.Get("post")
	// save id post in cookie
	object.CookiePostID(w, ck.PostString)
	idInt, err := strconv.Atoi(ck.PostString)
	if err != nil {
		api.Message(w, object.StatusByCodeAndLog(constant.Code400,
			err, "read handler: postID: atoi"))
		return
	}
	ck.Post = idInt
	// get data from db, parse and execute response
	hp.get(ctx, ck, w)
}

func (hp *hPost) CreatePost(ctx context.Context, ck *object.Cookie, sts object.Status,
	w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " ", r.URL.Path)
	ctx, cancel := context.WithTimeout(ctx, constant.TimeLimit)
	defer cancel()
	// check errors in cookie
	if sts != nil {
		api.Message(w, sts)
		return
	}
	if r.Method != "POST" {
		api.Message(w, object.StatusByCode(constant.Code405))
		return
	}
	// need session always to continue
	if !ck.Session {
		api.Message(w, object.StatusByText(constant.AccessDenied,
			"", nil))
		return
	}
	// create DTO with a new post
	post := dto.NewPost(nil, nil, ck)
	// add request data to DTO & check fields for valid
	if !post.Add(r) || !post.Valid() {
		api.Message(w, post.Obj.Sts)
		return
	}
	// create post in database
	id, sts := hp.service.Create(ctx, post)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	ck.Post = id
	// send new id to cookie
	object.CookiePostID(w, strconv.Itoa(id))
	// get data from db, parse and execute response
	hp.get(ctx, ck, w)
}

func (hp *hPost) CreateComment(ctx context.Context, ck *object.Cookie, sts object.Status,
	w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " ", r.URL.Path)
	ctx, cancel := context.WithTimeout(ctx, constant.TimeLimit)
	defer cancel()
	// check errors in cookie
	if sts != nil {
		api.Message(w, sts)
		return
	}
	if r.Method != "POST" {
		api.Message(w, object.StatusByCode(constant.Code405))
		return
	}
	// need session always to continue
	if !ck.Session {
		api.Message(w, object.StatusByText(constant.AccessDenied,
			"", nil))
		return
	}
	// create DTO with a new comment
	comment := dto.NewComment(nil, nil, ck)
	// add request data to DTO & check fields for valid
	if !comment.Add(r) || !comment.Valid() {
		api.Message(w, comment.Obj.Sts)
		return
	}
	// check postID valid
	idPost, sts := hp.service.Check(ctx, []string{ck.PostString})
	if sts != nil {
		api.Message(w, sts)
		return
	}
	if len(idPost) == 0 {
		api.Message(w, object.StatusByCode(constant.Code400))
		return
	}
	ck.Post = idPost[0]
	// create comment in database
	_, sts = hp.sComment.Create(ctx, comment)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// get data from db, parse and execute response
	hp.get(ctx, ck, w)
}

func (hp *hPost) CreateLike(ctx context.Context, ck *object.Cookie, sts object.Status,
	w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " ", r.URL.Path)
	ctx, cancel := context.WithTimeout(ctx, constant.TimeLimit)
	defer cancel()
	// check errors in cookie
	if sts != nil {
		api.Message(w, sts)
		return
	}
	if r.Method != "GET" {
		api.Message(w, object.StatusByCode(constant.Code405))
		return
	}
	// need session always to continue
	if !ck.Session {
		api.Message(w, object.StatusByText(constant.AccessDenied,
			"", nil))
		return
	}
	hp.ratio.Rate(ctx, ck, r)
	// get data from db, parse and execute response
	hp.get(ctx, ck, w)
}

func (hp *hPost) get(ctx context.Context, ck *object.Cookie, w http.ResponseWriter) {
	// parse
	pe, sts := api.NewParseExecute("post").Parse()
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// link for "form action" submit
	pe.Data["RatioLink"] = constant.URLReadRatio
	// insert session
	pe.Data["Session"] = ck.Session
	// create new model posts
	log.Printf("cookie is %v", ck)
	p := model.NewPosts(nil, ck)
	// make keys for sort posts
	p.MakeKeys(constant.KeyPost)
	// insert posts
	pe.Data["Posts"], sts = hp.service.Get(ctx, p)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// insert method to show - one post or all posts
	pe.Data["AllPost"] = p.St.AllPost
	// create new model categories
	c := model.NewCategories(nil, nil)
	// insert categories
	pe.Data["Category"], sts = hp.sCategory.GetList(ctx, c)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// create new model comments
	cm := model.NewComments(nil, ck)
	// make keys for sort posts
	cm.MakeKeys(constant.KeyPost)
	// insert comments
	pe.Data["Comments"], sts = hp.sComment.Get(ctx, cm)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// execute
	pe.Execute(w, constant.Code200)
}
