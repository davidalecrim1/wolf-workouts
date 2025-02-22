package handler

import (
	"context"

	gen "github.com/davidalecrim1/wolf-workouts/internal/trainer/handler/generated"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TrainerGrpcHandler struct {
	gen.UnimplementedTrainerServiceServer
	db *mongo.Client
}

func NewTrainerGrpcHandler(db *mongo.Client) *TrainerGrpcHandler {
	return &TrainerGrpcHandler{
		db: db,
	}
}

func (h *TrainerGrpcHandler) ScheduleTraining(ctx context.Context, in *gen.ScheduleHourRequest) (*emptypb.Empty, error) {
	panic("not implemented") // TODO: Implement
}

func (h *TrainerGrpcHandler) HealthCheck(ctx context.Context, in *emptypb.Empty) (*gen.HealthCheckResponse, error) {
	err := h.db.Ping(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Database connection failed")
	}

	return &gen.HealthCheckResponse{
		Message: "OK",
	}, nil
}
