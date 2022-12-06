package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/routers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"log"
	"net/http"
)

func main() {

	m := storage.CreateBaseStorage()

	r := routers.NewRouter(m)

	log.Println(http.ListenAndServe(":8080", r))
}
