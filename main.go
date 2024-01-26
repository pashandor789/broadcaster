package main

import (
	"context"
	telegram "github.com/pashandor789/broadcaster/bot"
	"github.com/pashandor789/broadcaster/http"
	"github.com/pashandor789/broadcaster/repository"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	err := SetupConfig()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := repository.NewPostgresSQLPool(GetDatabaseURL())
	if err != nil {
		log.Fatal(err)
	}

	tg, err := telegram.NewTgBot(GetBotConfig(), repo)

	go tg.Serve(context.Background())

	server, err := http.NewHTTPServer(GetServerConfig(), tg)
	if err != nil {
		log.Fatal(err)
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
