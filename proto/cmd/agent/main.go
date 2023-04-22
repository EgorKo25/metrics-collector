package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/EgorKo25/DevOps-Track-Yandex/proto/service"
)

func main() {

	var metric service.Metric

	metric.Id = "Random"
	metric.Type = service.Metric_GAUGE
	metric.Hash = ""
	metric.Value = 123.0

	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%s", err)
	}

	agent := service.NewServiceClient(conn)

	req, err := agent.TakeMetric(context.Background(), &service.MetricRequest{Metric: &metric})
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Println(req.Status)
}
