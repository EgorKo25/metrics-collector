package main

import (
	config "github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/reader"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/saver"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/serializer"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/handlers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/routers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"log"
	"net/http"
)

func main() {

	cfg := config.NewServerConfig()

	srl := serializer.NewSerialize()

	strg := storage.NewStorage()

	handler := handlers.NewHandler(strg, srl)

	router := routers.NewRouter(handler)

	save := saver.NewSave(cfg, strg, srl)

	read, _ := reader.NewRead(cfg, strg, srl)

	if cfg.Restore {
		read.ReadAll()
	}

	go save.Run()
	log.Println(http.ListenAndServe(cfg.Address, router))
}
