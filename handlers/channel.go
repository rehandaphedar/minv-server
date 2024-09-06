package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"git.sr.ht/~rehandaphedar/minv-server/db"
	"git.sr.ht/~rehandaphedar/minv-server/token"
	"git.sr.ht/~rehandaphedar/minv-server/validators"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type channelParams struct {
	Channelname string `json:"channelname" validate:"required,min=3,max=128"`
}

func SelectChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := db.Queries.ChannelSelectChannels(context.Background())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"error": err.Error(),
		})
		return
	}

	render.JSON(w, r, channels)
}

func SelectChannel(w http.ResponseWriter, r *http.Request) {
	var body channelParams
	body.Channelname = chi.URLParam(r, "channelname")

	err := validators.ValidateStruct(body)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"error": err.Error(),
		})
		return
	}

	channel, err := db.Queries.ChannelSelectChannel(context.Background(), body.Channelname)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Channel not found",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, channel)
}

func SelectLoggedInChannel(w http.ResponseWriter, r *http.Request) {
	tokenCookie, _ := r.Cookie("token")
	tokenString := tokenCookie.Value
	payload, _ := token.VerifyToken(tokenString)
	channelname := payload.Username

	channel, err := db.Queries.ChannelSelectChannel(context.Background(), channelname)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Unexpected error",
		})
		return
	}

	render.JSON(w, r, channel)
}

func DeleteChannel(w http.ResponseWriter, r *http.Request) {
	tokenCookie, _ := r.Cookie("token")
	tokenString := tokenCookie.Value
	payload, _ := token.VerifyToken(tokenString)
	channelname := payload.Username

	videos, err := db.Queries.VideoSelectVideosByChannel(context.Background(), channelname)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}

	channel, err := db.Queries.ChannelDeleteChannel(context.Background(), channelname)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error while deleting channel",
			"error":   err.Error(),
		})
		return
	}

	for _, video := range videos {
		os.Remove(fmt.Sprintf("./data/videos/%s.mp4", video.Slug))
	}

	render.JSON(w, r, channel)
}
