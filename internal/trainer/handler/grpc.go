package handler

import (
	"context"
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/trainer/app/command"
	gen "github.com/davidalecrim1/wolf-workouts/internal/trainer/handler/generated"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TrainerGrpcHandler struct {
	gen.UnimplementedTrainerServiceServer
	db                             *mongo.Client
	commandScheduleTrainingHandler *command.ScheduleTrainingHandler
}

func NewTrainerGrpcHandler(db *mongo.Client, cmdSth *command.ScheduleTrainingHandler) *TrainerGrpcHandler {
	return &TrainerGrpcHandler{
		commandScheduleTrainingHandler: cmdSth,
		db:                             db,
	}
}

func (h *TrainerGrpcHandler) ScheduleTraining(ctx context.Context, in *gen.ScheduleHourRequest) (*emptypb.Empty, error) {
	timeStr := in.GetTime()
	if timeStr == "" {
		zap.S().Error("Invalid time sent to ScheduleTraining", "time", timeStr, "ctx", ctx)
		return nil, status.Errorf(codes.InvalidArgument, "Time is required")
	}

	trainingTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		zap.S().Error("Invalid time parsing to ScheduleTraining", "error", err, "ctx", ctx)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid time format: %v", err)
	}

	cmd := &command.ScheduleTrainingCommand{
		Timestamp: trainingTime,
	}

	err = h.commandScheduleTrainingHandler.Handle(ctx, cmd)
	if err != nil {
		zap.S().Error("Failed to ScheduleTraining", "error", err, "ctx", ctx)
		return nil, status.Errorf(codes.Internal, "Failed to ScheduleTraining")
	}

	return &emptypb.Empty{}, nil
}

func (h *TrainerGrpcHandler) HealthCheck(ctx context.Context, in *emptypb.Empty) (*gen.HealthCheckResponse, error) {
	err := h.db.Ping(ctx, nil)
	if err != nil {
		zap.S().Error("Failed to HealthCheck", "error", err, "ctx", ctx)
		return nil, status.Errorf(codes.Internal, "Database connection failed")
	}

	return &gen.HealthCheckResponse{
		Message: "OK",
	}, nil
}
