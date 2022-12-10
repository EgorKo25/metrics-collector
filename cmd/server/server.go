package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/routers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"log"
	"net/http"
)

func main() {

	memStorage := storage.NewStorage()

	router := routers.NewRouter(memStorage)

	log.Println(http.ListenAndServe(":8080", router))
}
