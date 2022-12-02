package main

import (
	"DevOps-Track-Yandex/internal/ServerSupport/routers"
	"DevOps-Track-Yandex/internal/StorageSupport"
	"log"
	"net/http"
)

func main() {
	var m StorageSupport.MemStats
	m.CreateBaseMap()

	r := routers.NewRouter(m)

	log.Fatal(http.ListenAndServe(":8080", r))
}
