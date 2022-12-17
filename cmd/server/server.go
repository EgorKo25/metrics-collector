package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/handlers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/routers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"log"
	"net/http"
)

func main() {

	strg := storage.NewStorage()

	handler := handlers.NewHandler(strg)

	router := routers.NewRouter(handler)

	log.Println(http.ListenAndServe(":8080", router))
}
