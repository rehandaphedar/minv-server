package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.sr.ht/~rehandaphedar/minv-server/db"
	"git.sr.ht/~rehandaphedar/minv-server/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/feeds"
	"github.com/spf13/viper"
)

func AtomSelectVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := db.Queries.VideoSelectVideos(context.Background())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"error": err.Error(),
		})
		return
	}

	feed := &feeds.Feed{
		Title:       "minv videos",
		Link:        &feeds.Link{Href: fmt.Sprintf("%s/static", viper.GetString("base_url"))},
		Description: "Main minv videos feed",
		Author:      &feeds.Author{Name: "minv instance users"},
		Created:     time.Now(),
	}

	feed.Items = []*feeds.Item{}

	for _, video := range videos {
		created, err := time.Parse(viper.GetString("db_time_format"), video.Uploaded)
		if err != nil {
			log.Printf("Error while formatting string: %v", err)
			continue
		}

		feed.Items = append(feed.Items, &feeds.Item{
			Title:       video.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/static/%s.mp4", viper.GetString("base_url"), video.Slug)},
			Description: video.Description,
			Author:      &feeds.Author{Name: video.Uploader},
			Created:     created,
		})
	}

	util.RenderAtom(w, r, feed)
}

func AtomSelectChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := db.Queries.ChannelSelectChannels(context.Background())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	feed := &feeds.Feed{
		Title:       "minv channels",
		Link:        &feeds.Link{Href: fmt.Sprintf("%s/channels", viper.GetString("base_url"))},
		Description: "Main minv channels feed",
		Author:      &feeds.Author{Name: "minv instance users"},
		Created:     time.Now(),
	}

	feed.Items = []*feeds.Item{}

	for _, channel := range channels {
		created, err := time.Parse(viper.GetString("db_time_format"), channel.Created)
		if err != nil {
			log.Printf("Error while formatting string: %v", err)
			continue
		}

		feed.Items = append(feed.Items, &feeds.Item{
			Title:       channel.Channelname,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/channel/%s", viper.GetString("base_url"), channel.Channelname)},
			Description: fmt.Sprintf("%s's  channel", channel.Channelname),
			Author:      &feeds.Author{Name: channel.Channelname},
			Created:     created,
		})
	}

	util.RenderAtom(w, r, feed)
}

func AtomSelectVideosByChannel(w http.ResponseWriter, r *http.Request) {
	channelname := chi.URLParam(r, "channelname")

	videos, err := db.Queries.VideoSelectVideosByChannel(context.Background(), channelname)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}

	feed := &feeds.Feed{
		Title:       fmt.Sprintf("%s's videos", channelname),
		Link:        &feeds.Link{Href: fmt.Sprintf("%s/channel/%s", viper.GetString("base_url"), channelname)},
		Description: fmt.Sprintf("Videos uploaded by %s", channelname),
		Author:      &feeds.Author{Name: channelname},
		Created:     time.Now(),
	}

	feed.Items = []*feeds.Item{}

	for _, video := range videos {
		created, err := time.Parse(viper.GetString("db_time_format"), video.Uploaded)
		if err != nil {
			log.Printf("Error while formatting string: %v", err)
			continue
		}

		feed.Items = append(feed.Items, &feeds.Item{
			Title:       video.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/static/%s.mp4", viper.GetString("base_url"), video.Slug)},
			Description: video.Description,
			Author:      &feeds.Author{Name: video.Uploader},
			Created:     created,
		})
	}

	util.RenderAtom(w, r, feed)
}
