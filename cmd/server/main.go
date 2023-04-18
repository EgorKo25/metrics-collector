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
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/database"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/encryption"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/file"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/hashing"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/handlers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/middleware"
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

	cpr := middleware.NewCompressor()

	enc, err := encryption.NewEncryptor(cfg.CryptoKey, "private")
	if err != nil {
		log.Fatalf("%s", err)
	}

	handler := handlers.NewHandler(str, cpr, hsr, db, enc)

	middle := middleware.NewMiddle(cfg)

	router := routers.NewRouter(handler, middle)

	save := file.NewSave(cfg, str)

	go func() {
		err = save.Run()
		if err != nil {
			log.Println("save file error: ", err)
		}
	}()

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err = server.Shutdown(context.Background()); err != nil {

			log.Printf("\nHTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-idleConnsClosed
}
