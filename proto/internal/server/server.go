package server

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/database"
	"github.com/EgorKo25/DevOps-Track-Yandex/proto/service"
)

type Server struct {
	db *database.DB
}

func NewServer(db *database.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) TakeMetric(ctx context.Context, in *service.MetricRequest) (*service.MetricResponse, error) {

	childCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if in.Metric.Type.Number() == 1 {
		value, err := strconv.ParseFloat(in.Metric.Value.String(), 64)
		if err != nil {
			return &service.MetricResponse{
				Metric: nil,
				Status: int32(codes.InvalidArgument),
			}, err
		}
		_, _ = s.db.DB.ExecContext(childCtx, "INSERT INTO metrics (id, type, hash, value) VALUES ($1, $2, $3, $4)",
			in.Metric.Id, in.Metric.Type, in.Metric.Hash, value,
		)
	}
	if in.Metric.Type.Number() == 2 {
		delta, err := strconv.Atoi(in.Metric.Value.String())
		if err != nil {
			return &service.MetricResponse{
				Metric: nil,
				Status: int32(codes.InvalidArgument),
			}, err
		}
		_, _ = s.db.DB.ExecContext(childCtx, "INSERT INTO metrics (id, type, hash, delta) VALUES ($1, $2, $3, $4)",
			in.Metric.Id, in.Metric.Type.String(), in.Metric.Hash, delta,
		)
	}

	return &service.MetricResponse{
		Metric: in.Metric,
		Status: int32(codes.OK),
	}, nil
}

func (*Server) GetMetric(ctx context.Context, in *service.MetricRequest) (*service.MetricResponse, error) {
	return nil, nil
}

func (*Server) Ping(ctx context.Context, in *service.PingRequest) (*service.PingResponse, error) {
	return nil, nil
}

func (*Server) AllStats(ctx context.Context, in *service.AllStatsRequest) (*service.AllStatsResponse, error) {
	return nil, nil
}
