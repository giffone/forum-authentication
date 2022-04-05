package api

import (
	"github.com/giffone/forum-authentication/internal/object"
	"log"
	"net/http"
)

func Message(w http.ResponseWriter, s object.Status) {
	pe, status := NewParseExecute("message").Parse()
	if status != nil {
		sts := status.Status()
		http.Error(w, sts.StatusBody, sts.StatusCode)
		return
	}
	sts := s.Status()
	if sts.StatusCode < 200 && sts.StatusCode > 500 {
		pe.Data["StatusCode"] = "Message"
		w.WriteHeader(http.StatusNoContent)
		log.Printf("message: writing code to w: %d", http.StatusNoContent)
	} else {
		pe.Data["StatusCode"] = http.StatusText(sts.StatusCode)
		w.WriteHeader(sts.StatusCode)
		log.Printf("message: writing code to w: %d", sts.StatusCode)
	}
	pe.Data["Return"] = sts.ReturnPage
	pe.Data["StatusBody"] = sts.StatusBody

	pe.Execute(w, 0)
}
