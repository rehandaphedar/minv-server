package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"

	"git.sr.ht/~rehandaphedar/minv-server/config"
	"git.sr.ht/~rehandaphedar/minv-server/db"
	"git.sr.ht/~rehandaphedar/minv-server/handlers"
	"git.sr.ht/~rehandaphedar/minv-server/routes"
	"git.sr.ht/~rehandaphedar/minv-server/validators"
)

func main() {
	os.Mkdir("./data", os.ModePerm)
	os.Mkdir("./data/temp-upload", os.ModePerm)
	os.Mkdir("./data/videos", os.ModePerm)

	config.InitialiseConfig("./data")
	db.Connect()
	validators.Initialise()

	r := chi.NewRouter()
	routes.Setup(r)

	workDir, _ := os.Getwd()
	staticDir := http.Dir(filepath.Join(workDir, "data", "videos"))
	handlers.FileServer(r, "/static", staticDir)

	http.ListenAndServe(fmt.Sprintf(":%s", viper.GetString("port")), r)
}
