package main

import (
	"log"
	"net/http"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/middleware"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/reader"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/saver"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/serializer"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/handlers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/routers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
)

func main() {

	cfg := config.NewServerConfig()

	srl := serializer.NewSerialize()

	storage := storage.NewStorage()

	compressor := middleware.NewCompressor()

	handler := handlers.NewHandler(storage, srl, compressor)

	router := routers.NewRouter(handler)

	save := saver.NewSave(cfg, storage, srl)

	read, _ := reader.NewRead(cfg, storage, srl)

	if cfg.Restore {
		read.ReadAll()
	}

	go save.Run()
	log.Println(http.ListenAndServe(cfg.Address, router))
}
