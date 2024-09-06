package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"git.sr.ht/~rehandaphedar/minv-server/db"
	"git.sr.ht/~rehandaphedar/minv-server/token"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
	"github.com/ulule/deepcopier"

	"git.sr.ht/~rehandaphedar/minv-server/sqlc"
	"git.sr.ht/~rehandaphedar/minv-server/validators"
	"golang.org/x/crypto/bcrypt"
)

type authParams struct {
	Channelname string `json:"channelname" validate:"required,min=3,max=128"`
	Password    string `json:"password" validate:"required,min=3,max=128"`
}

type authResponse struct {
	Channelname string `json:"channelname"`
	Created     string `json:"created"`
}

func Auth(w http.ResponseWriter, r *http.Request) {

	var body authParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"error": "Could not parse request body",
		})
		return
	}

	err := validators.ValidateStruct(body)
	if err != nil {

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"error": err.Error(),
		})

		return
	}

	channel, err := db.Queries.AuthSelectChannel(context.Background(), body.Channelname)

	if err != nil { // Implies that the channel doesn't exist
		register(w, r, body)
	} else {
		login(w, r, body, channel)
	}
}

func register(w http.ResponseWriter, r *http.Request, body authParams) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	registerParams := sqlc.AuthInsertChannelParams{
		Channelname: body.Channelname,
		Password:    string(hashedPassword),
	}

	channel, err := db.Queries.AuthInsertChannel(context.Background(), registerParams)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	var response authResponse
	deepcopier.Copy(&channel).To(&response)

	addTokenCookie(w, r, response.Channelname)
	render.JSON(w, r, response)
}

func login(w http.ResponseWriter, r *http.Request, body authParams, channel sqlc.Channel) {

	err := bcrypt.CompareHashAndPassword([]byte(channel.Password), []byte(body.Password))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"error": "Invalid channelname or password",
		})
		return
	}

	var response authResponse
	deepcopier.Copy(&channel).To(&response)

	addTokenCookie(w, r, response.Channelname)
	render.JSON(w, r, response)
}

func addTokenCookie(w http.ResponseWriter, r *http.Request, channelname string) {

	pasetoDuration := viper.GetDuration("paseto_duration")

	token, err := token.CreateToken(channelname, pasetoDuration)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}
	tokenCookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(pasetoDuration),
	}
	http.SetCookie(w, &tokenCookie)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(-time.Hour),
	})
}
