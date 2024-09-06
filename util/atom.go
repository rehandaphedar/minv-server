package util

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/gorilla/feeds"
)

func RenderAtom(w http.ResponseWriter, r *http.Request, feed *feeds.Feed) {
    atom, err := feed.ToAtom()
    if err != nil {
		render.Status(r, http.StatusInternalServerError)
		log.Println(err)
		return
    }

	w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")
	w.Write([]byte(atom))
}
