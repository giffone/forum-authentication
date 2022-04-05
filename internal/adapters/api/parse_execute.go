package api

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"html/template"
	"log"
	"net/http"
	"time"
)

type ParseExecute struct {
	PathTmpl   []string
	DefineTmpl string
	Data       map[string]interface{}
	tmpl       *template.Template
}

func NewParseExecute(define string) *ParseExecute {
	switch define {
	case "index":
		return &ParseExecute{
			PathTmpl: []string{
				constant.PathIndex, constant.PathHeaderObj,
				constant.PathFooterObj, constant.PathPostObj,
				constant.PathCategoriesObj, constant.PathLikesObj,
				constant.PathPostObj,
			},
			DefineTmpl: define,
			Data:       make(map[string]interface{}),
		}
	case "post":
		return &ParseExecute{
			PathTmpl: []string{
				constant.PathPost, constant.PathHeaderObj,
				constant.PathPostObj, constant.PathCategoriesObj,
				constant.PathCommentsObj, constant.PathLikesObj,
				constant.PathFooterObj,
			},
			DefineTmpl: define,
			Data:       make(map[string]interface{}),
		}
	case "account":
		return &ParseExecute{
			PathTmpl: []string{
				constant.PathAccount, constant.PathAccountUser,
				constant.PathHeaderObj, constant.PathPostObj,
				constant.PathCategoriesObj, constant.PathCommentsObj,
				constant.PathLikesObj, constant.PathFooterObj,
			},
			DefineTmpl: define,
			Data:       make(map[string]interface{}),
		}
	case "login":
		return &ParseExecute{
			PathTmpl:   []string{constant.PathLoginObj},
			DefineTmpl: define,
			Data:       nil,
		}
	case "message":
		return &ParseExecute{
			PathTmpl:   []string{constant.PathMessage},
			DefineTmpl: define,
			Data:       make(map[string]interface{}),
		}
	default:
		return &ParseExecute{
			DefineTmpl: "nil",
		}
	}
}

func (pe *ParseExecute) Parse() (*ParseExecute, object.Status) {
	if pe.DefineTmpl == "nil" {
		return nil, object.StatusByCodeAndLog(constant.Code500,
			nil, "parseFile: unknown define")
	}
	var err error
	myFunc := template.FuncMap{
		"dateForum": func(t time.Time) string {
			return t.Format(constant.ForumLayoutDate)
		},
	}
	pe.tmpl, err = template.New("").Funcs(myFunc).ParseFiles(pe.PathTmpl...)
	if err != nil {
		return nil, object.StatusByCodeAndLog(constant.Code500,
			err, "parseFile:")
	}
	return pe, nil
}

func (pe *ParseExecute) Execute(w http.ResponseWriter, code int) {
	if code != 0 {
		w.WriteHeader(code)
		log.Printf("execute: writing code to w: %d", code)
	}

	err := pe.tmpl.ExecuteTemplate(w, pe.DefineTmpl, pe.Data)
	//log.Printf("execute: w header is: %v", w)
	if err != nil {
		sts := object.StatusByCodeAndLog(constant.Code500,
			err, "executeTemplate:")
		Message(w, sts)
	}
}
