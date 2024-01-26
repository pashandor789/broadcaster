package http

import (
	"github.com/gorilla/mux"
	telegram "github.com/pashandor789/broadcaster/bot"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func NewHTTPServer(cfg ServerConfig, bot *telegram.TgBot) (*http.Server, error) {
	r := mux.NewRouter()

	h := broadcastHandler{TgBot: bot}
	r.Methods(http.MethodPost).PathPrefix("/broadcast").Handler(h)

	return &http.Server{
		Addr:         "localhost:" + strconv.Itoa(int(cfg.Port)),
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}, nil
}

type broadcastHandler struct {
	*telegram.TgBot
}

func (b broadcastHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = b.TgBot.BroadcastSubscribers(req.Context(), string(content))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
