package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/database"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"log"
	"net"

	"google.golang.org/grpc"

	config "github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	server "github.com/EgorKo25/DevOps-Track-Yandex/proto/internal/server"
	service "github.com/EgorKo25/DevOps-Track-Yandex/proto/service"
)

func main() {

	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("%s", err)
	}

	str := storage.NewStorage()

	db := database.NewDB(cfg, str)

	srv := server.NewServer(db)
	s := grpc.NewServer()

	service.RegisterServiceServer(s, srv)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("%s", err)
	}

	if err := s.Serve(l); err != nil {
		log.Fatalf("%s", err)
	}
}
