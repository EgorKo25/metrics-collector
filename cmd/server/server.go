package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/serializer"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/handlers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/routers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"log"
	"net/http"
)

func main() {

	srl := serializer.NewSerialize()

	strg := storage.NewStorage()

	handler := handlers.NewHandler(strg, srl)

	router := routers.NewRouter(handler)

	log.Println(http.ListenAndServe(":8080", router))
}
