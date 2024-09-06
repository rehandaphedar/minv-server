package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"git.sr.ht/~rehandaphedar/minv-server/db"
	"git.sr.ht/~rehandaphedar/minv-server/sqlc"
	"git.sr.ht/~rehandaphedar/minv-server/token"
	"git.sr.ht/~rehandaphedar/minv-server/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gosimple/slug"
	"github.com/spf13/viper"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func SelectVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := db.Queries.VideoSelectVideos(context.Background())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"error": err.Error(),
		})
		return
	}

	render.JSON(w, r, videos)
}

func SelectVideo(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	video, err := db.Queries.VideoSelectVideo(context.Background(), slug)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Video not found",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, video)
}

func SelectVideosByChannel(w http.ResponseWriter, r *http.Request) {
	channelname := chi.URLParam(r, "channelname")

	videos, err := db.Queries.VideoSelectVideosByChannel(context.Background(), channelname)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Channel not found",
			"error":   err.Error(),
		})
	}

	render.JSON(w, r, videos)
}

func CreateVideo(w http.ResponseWriter, r *http.Request) {
	// 10 << 24 specifies a maximum upload of 10000 MB files
	r.ParseMultipartForm(10 << 24)

	// Get user info
	var videoData sqlc.Video

	tokenCookie, _ := r.Cookie("token")
	tokenString := tokenCookie.Value

	payload, _ := token.VerifyToken(tokenString)
	videoData.Uploader = payload.Username

	// Retrieve the file
	file, handler, err := r.FormFile("minv--video")
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error retrieving the file",
			"error":   err.Error(),
		})
		return
	}

	defer file.Close()

	// Ensure that the file is a video
	if fileType := handler.Header.Get("Content-Type"); !util.IsVideoFile(fileType) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "The file is not a video",
		})
		return
	}

	// Get video data
	err = json.Unmarshal([]byte(r.FormValue("minv--data")), &videoData)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Invalid data",
			"error":   err.Error(),
		})
		return
	}

	videoData.Slug = slug.Make(videoData.Title)

	slugSuffix := 0

	insertVideoParams := &sqlc.VideoInsertVideoParams{
		Slug:        videoData.Slug,
		Title:       videoData.Title,
		Description: videoData.Description,
		Uploader:    videoData.Uploader,
	}
	returnedVideo, err := db.Queries.VideoInsertVideo(context.Background(), *insertVideoParams)
	for ; err != nil; returnedVideo, err = db.Queries.VideoInsertVideo(context.Background(), *insertVideoParams) {
		slugSuffix += 1
		insertVideoParams.Slug = fmt.Sprintf("%s:%d", videoData.Slug, slugSuffix)
	}

	go processVideo(returnedVideo, file)
	render.JSON(w, r, returnedVideo)
}

func processVideo(returnedVideo sqlc.Video, file multipart.File) {
	// Save the file

	// Create a temporary file within our temp-upload directory that follows a particular naming pattern
	tempFile, err := os.CreateTemp("data/temp-upload", "upload-*")
	if err != nil {
		db.Queries.VideoUpdateVideo(context.Background(), sqlc.VideoUpdateVideoParams{
			Slug:        returnedVideo.Slug,
			Title:       returnedVideo.Title,
			Description: returnedVideo.Description,
			Processed:   2,
		})
		return
	}
	defer tempFile.Close()

	// Read all of the contents of our uploaded file into a byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		db.Queries.VideoUpdateVideo(context.Background(), sqlc.VideoUpdateVideoParams{
			Slug:        returnedVideo.Slug,
			Title:       returnedVideo.Title,
			Description: returnedVideo.Description,
			Processed:   2,
		})
		return
	}
	// Write this byte array to our temporary file
	tempFile.Write(fileBytes)

	resolution := viper.GetIntSlice("resolution")
	x, y := resolution[0], resolution[1]

	err = ffmpeg.Input(tempFile.Name()).
		Output(fmt.Sprintf("./data/videos/%s.mp4", returnedVideo.Slug), ffmpeg.KwArgs{
			"c:v": "libx264",
			"vf":  fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2,setsar=1", x, y, x, y)}).
		OverWriteOutput().Run()

	if err != nil {
		db.Queries.VideoUpdateVideo(context.Background(), sqlc.VideoUpdateVideoParams{
			Slug:        returnedVideo.Slug,
			Title:       returnedVideo.Title,
			Description: returnedVideo.Description,
			Processed:   2,
		})
		return
	}

	// Delete the temporary file and return the information
	os.Remove(tempFile.Name())

	// Video was processed successfully
	db.Queries.VideoUpdateVideo(context.Background(), sqlc.VideoUpdateVideoParams{
		Slug:        returnedVideo.Slug,
		Title:       returnedVideo.Title,
		Description: returnedVideo.Description,
		Processed:   1,
	})
}

func DeleteVideo(w http.ResponseWriter, r *http.Request) {
	tokenCookie, _ := r.Cookie("token")
	tokenString := tokenCookie.Value
	payload, _ := token.VerifyToken(tokenString)
	channelname := payload.Username

	slug := chi.URLParam(r, "slug")

	returnedVideo, err := db.Queries.VideoSelectVideo(context.Background(), slug)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Video doesn't exist",
			"error":   err.Error(),
		})
		return
	}

	if returnedVideo.Uploader != channelname {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{
			"message": "You do not own this video",
			"error":   err.Error(),
		})
		return
	}

	returnedVideo, err = db.Queries.VideoDeleteVideo(context.Background(), slug)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error while deleting video",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, returnedVideo)
	os.Remove(fmt.Sprintf("./data/videos/%s.mp4", returnedVideo.Slug))
}

func UpdateVideo(w http.ResponseWriter, r *http.Request) {
	tokenCookie, _ := r.Cookie("token")
	tokenString := tokenCookie.Value
	payload, _ := token.VerifyToken(tokenString)
	channelname := payload.Username

	slug := chi.URLParam(r, "slug")

	returnedVideo, err := db.Queries.VideoSelectVideo(context.Background(), slug)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Video doesn't exist",
			"error":   err.Error(),
		})
		return
	}

	if returnedVideo.Uploader != channelname {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{
			"message": "You do not own this video",
			"error":   err.Error(),
		})
		return
	}

	updatedVideoData := &sqlc.VideoUpdateVideoParams{
		Title:       returnedVideo.Title,
		Description: returnedVideo.Description,
		Slug:        returnedVideo.Slug,
		Processed:   returnedVideo.Processed,
	}
	var body sqlc.VideoUpdateVideoParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"error": "Could not parse request body",
		})
		return
	}

	if body.Title != "" {
		updatedVideoData.Title = body.Title
	}
	if body.Description != "" {
		updatedVideoData.Description = body.Description
	}

	updatedVideo, err := db.Queries.VideoUpdateVideo(context.Background(), *updatedVideoData)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Unexpected error",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, updatedVideo)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
