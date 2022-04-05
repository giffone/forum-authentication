package account

import (
	"context"
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/model"
	"github.com/giffone/forum-authentication/internal/service"
	"log"
	"net/http"
)

type hAccount struct {
	sPost     service.Post
	sCategory service.Category
	sComment  service.Comment
	ratio     api.Ratio
}

func NewHandler(sPost service.Post, sCategory service.Category,
	sComment service.Comment, ratio api.Ratio) api.Handler {
	return &hAccount{
		sPost:     sPost,
		sCategory: sCategory,
		sComment:  sComment,
		ratio:     ratio,
	}
}

func (ha *hAccount) Register(ctx context.Context, router *http.ServeMux, session api.Session) {
	router.HandleFunc(constant.URLAccount, session.Check(ctx, ha.ByUser))
	router.HandleFunc(constant.URLAccountRatio, session.Check(ctx, ha.CreateLike))
}

func (ha *hAccount) ByUser(ctx context.Context, ck *object.Cookie, sts object.Status,
	w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " ", r.URL.Path)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	if r.Method != "GET" {
		api.Message(w, object.StatusByCode(constant.Code405))
		return
	}
	// clear unused postID
	object.CookiePostIDDel(w)
	// get data from db, parse and execute response
	ha.get(ctx, ck, w)
}

func (ha *hAccount) CreateLike(ctx context.Context, ck *object.Cookie, sts object.Status,
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
	// make ratio
	ha.ratio.Rate(ctx, ck, r)
	redirect := r.PostFormValue(constant.KeyLink)
	http.Redirect(w, r, constant.URLAccount+redirect, constant.Code302)
	//ha.get(ctx, ck, w)
}

func (ha *hAccount) get(ctx context.Context, ck *object.Cookie, w http.ResponseWriter) {
	ctx, cancel := context.WithTimeout(ctx, constant.TimeLimit)
	defer cancel()
	// parse
	pe, sts := api.NewParseExecute("account").Parse()
	if sts != nil {
		api.Message(w, sts)
		return
	}
	//
	pe.Data["Acc"] = true
	// link for "form action" submit
	pe.Data["RatioLink"] = constant.URLAccountRatio
	// insert session
	pe.Data["Session"] = ck.Session
	// create new model posts
	posts := model.NewPosts(nil, ck)
	// make keys for sort posts
	posts.MakeKeys(constant.KeyUser)
	// insert posts
	pe.Data["Posts"], sts = ha.sPost.Get(ctx, posts)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// create new model posts
	postsRated := model.NewPosts(nil, ck)
	// make keys for sort posts
	postsRated.MakeKeys(constant.KeyRated)
	// insert posts
	pe.Data["PostsRated"], sts = ha.sPost.Get(ctx, postsRated)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// need make refer to post
	st := object.NewSettings()
	st.Refers = true
	// create new model comments
	comments := model.NewComments(st, ck)
	// make keys for sort posts
	comments.MakeKeys(constant.KeyUser)
	// insert posts
	pe.Data["Comments"], sts = ha.sComment.Get(ctx, comments)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// create new model comments
	commentsRated := model.NewComments(st.ClearKey(), ck)
	// make keys for sort posts
	commentsRated.MakeKeys(constant.KeyRated)
	// insert posts
	pe.Data["CommentsRated"], sts = ha.sComment.Get(ctx, commentsRated)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// create new model categories
	c := model.NewCategories(nil, nil)
	// insert categories
	pe.Data["Category"], sts = ha.sCategory.GetList(ctx, c)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// insert method to show - one post or all posts
	pe.Data["AllPost"] = posts.St.AllPost
	// execute
	pe.Execute(w, constant.Code200)
}
