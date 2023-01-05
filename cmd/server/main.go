package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/file"
	"log"
	"net/http"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/middleware"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/serializer"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/handlers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/routers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
)

func main() {

	cfg := config.NewServerConfig()

	srl := serializer.NewSerialize()

	str := storage.NewStorage()

	compressor := middleware.NewCompressor()

	handler := handlers.NewHandler(str, srl, compressor)

	router := routers.NewRouter(handler)

	save := file.NewSave(cfg, str)

	read, _ := file.NewRead(cfg, str)

	if cfg.Restore {
		read.ReadAll()
	}

	go save.Run()
	log.Println(http.ListenAndServe(cfg.Address, router))
}
