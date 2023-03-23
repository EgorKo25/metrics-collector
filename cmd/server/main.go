// Пакет с логикой сервера
//
// Приложение собирается командой:
//
//	go build server
//
// Запускается:
//
//	./agent
//
// Или
//
//	go run main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/database"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/file"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/hashing"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/middleware"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/handlers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/routers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {

	fmt.Printf("Build version: %s\nBuild date:    %s\nBuild commit:  %s\n", buildVersion, buildDate, buildCommit)

	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("%s: %s", config.ErrFlagParse, err)
	}

	str := storage.NewStorage()

	_, err = file.NewRead(cfg, str)
	if err != nil {
		log.Println("file read error: ", err)
	}

	hsr := hashing.NewHash(cfg.Key)

	db := database.NewDB(cfg, str)

	compressor := middleware.NewCompressor()

	handler := handlers.NewHandler(str, compressor, hsr, db)

	router := routers.NewRouter(handler)

	save := file.NewSave(cfg, str)

	go func() {
		err = save.Run()
		if err != nil {
			log.Println("save file error: ", err)
		}
	}()

	log.Println(http.ListenAndServe(cfg.Address, router))

}
