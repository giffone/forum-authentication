package dto

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"net/http"
	"time"
)

const (
	like    = 1
	dislike = 2
)

type Ratio struct {
	PostOrComm map[string]string // post or comment - map[post]id
	Ratio      int               // like/dislike
	Obj        object.Obj
}

func NewRatio(st *object.Settings, sts *object.Statuses, ck *object.Cookie) *Ratio {
	l := new(Ratio)
	l.Obj.NewObjects(st, sts, ck)
	return l
}

func (r *Ratio) AddByPOST(request *http.Request) bool {
	// ratio - like/dislike
	rate := request.PostFormValue(constant.KeyRate)
	if rate == constant.KeyLike {
		r.Ratio = like
	} else if rate == constant.KeyDislike {
		r.Ratio = dislike
	} else {
		r.Obj.Sts.StatusByCodeAndLog(constant.Code400,
			nil, "dto: like: postFormValue is nil")
		return false
	}
	// post or comment ID
	postOrComm := constant.KeyPost
	objID := request.PostFormValue(postOrComm)
	if objID == "" {
		postOrComm = constant.KeyComment
		objID = request.PostFormValue(postOrComm)
	}
	if objID == "" {
		r.Obj.Sts.StatusByCodeAndLog(constant.Code400,
			nil, "dto: like: post or comment id is nil")
		return false
	}
	r.PostOrComm = make(map[string]string)
	r.PostOrComm[postOrComm] = objID
	return true
}

func (r *Ratio) AddByGET(request *http.Request) bool {
	// get post id
	sts := r.Obj.Ck.CookiePostIDRead(request)
	if sts != nil {
		r.Obj.Sts = sts.Status()
		return false
	}
	// read url
	u := request.URL.Query()
	// ratio - like/dislike
	rate := u.Get(constant.KeyRate)
	if rate == constant.KeyLike {
		r.Ratio = like
	} else if rate == constant.KeyDislike {
		r.Ratio = dislike
	} else {
		r.Obj.Sts.StatusByCodeAndLog(constant.Code400,
			nil, "dto: like: postFormValue is nil")
		return false
	}
	// post or comment ID
	postOrComm := constant.KeyPost
	objID := u.Get(postOrComm)
	if objID == "" {
		postOrComm = constant.KeyComment
		objID = u.Get(postOrComm)
	}
	if objID == "" {
		r.Obj.Sts.StatusByCodeAndLog(constant.Code400,
			nil, "dto: like: post or comment id is nil")
		return false
	}
	r.PostOrComm = make(map[string]string)
	r.PostOrComm[postOrComm] = objID
	return true
}

func (r *Ratio) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		r.Obj.St.Key[key] = data
	} else {
		r.Obj.St.Key[constant.KeyPost] = []interface{}{0}
	}
}

// Create prepares query for db and choose fields for adding incoming data
func (r *Ratio) Create() *object.QuerySettings {
	qs := new(object.QuerySettings)
	qs.QueryName = constant.QueInsert4
	if _, ok := r.Obj.St.Key[constant.KeyPost]; ok {
		qs.QueryFields = []interface{}{
			constant.TabPostsLikes,
			constant.FieldUser,
			constant.FieldPost,
			constant.FieldLike,
			constant.FieldCreated,
		}
		qs.Fields = []interface{}{
			r.Obj.Ck.User,
			r.PostOrComm[constant.KeyPost],
			r.Ratio,
			time.Now(),
		}
	} else {
		qs.QueryFields = []interface{}{
			constant.TabCommentsLikes,
			constant.FieldUser,
			constant.FieldComment,
			constant.FieldLike,
			constant.FieldCreated,
		}
		qs.Fields = []interface{}{
			r.Obj.Ck.User,
			r.PostOrComm[constant.KeyComment],
			r.Ratio,
			time.Now(),
		}
	}
	return qs
}

func (r *Ratio) Delete() *object.QuerySettings {
	qs := new(object.QuerySettings)
	qs.QueryName = constant.QueDeleteBy
	if value, ok := r.Obj.St.Key[constant.KeyPost]; ok {
		qs.QueryFields = []interface{}{
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.FieldID,
		}
		qs.Fields = value
	} else if value, ok := r.Obj.St.Key[constant.KeyComment]; ok {
		qs.QueryFields = []interface{}{
			constant.TabCommentsLikes,
			constant.TabCommentsLikes,
			constant.FieldID,
		}
		qs.Fields = value // for null list
	} else {
		qs.QueryFields = []interface{}{
			constant.TabCommentsLikes,
			constant.TabCommentsLikes,
			constant.FieldID,
		}
		qs.Fields = []interface{}{0} // for null list
	}
	return qs
}
