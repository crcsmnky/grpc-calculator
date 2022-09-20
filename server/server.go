package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/crcsmnky/grpc-calculator/config"
	pb "github.com/crcsmnky/grpc-calculator/proto"

	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
)

var _ pb.CalculatorServer = (*Server)(nil)

type Server struct {
	pb.UnimplementedCalculatorServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Calculate(ctx context.Context, r *pb.BinaryOperation) (*pb.CalculationResult, error) {
	log.Printf("[server:Calculate] Started")
	if ctx.Err() == context.Canceled {
		return &pb.CalculationResult{}, fmt.Errorf("client cancelled: abandoning")
	}

	switch r.GetOperation() {			
		case pb.Operation_ADD:
			return add(ctx, r)
		case pb.Operation_SUBTRACT:
			return subtract(ctx, r)
		case pb.Operation_MULTIPLY:
			time.Sleep(125 * time.Millisecond)
			return multiply(ctx, r)
		case pb.Operation_DIVIDE:
			time.Sleep(250 * time.Millisecond)
			return divide(ctx, r)
		default:
			return &pb.CalculationResult{}, fmt.Errorf("undefined operation")
	}
}

func add(ctx context.Context, r *pb.BinaryOperation) (*pb.CalculationResult, error) {
	tracer := otel.GetTracerProvider().Tracer("tracer")
	_, span := tracer.Start(ctx, "add")
	defer span.End()

	return &pb.CalculationResult{Result: r.GetFirstOperand() + r.GetSecondOperand()}, nil
}

func subtract(ctx context.Context, r *pb.BinaryOperation) (*pb.CalculationResult, error) {
	tracer := otel.GetTracerProvider().Tracer("tracer")
	_, span := tracer.Start(ctx, "subtract")
	defer span.End()

	return &pb.CalculationResult{Result: r.GetFirstOperand() - r.GetSecondOperand()}, nil	
}

func multiply(ctx context.Context, r *pb.BinaryOperation) (*pb.CalculationResult, error) {
	tracer := otel.GetTracerProvider().Tracer("tracer")
	_, span := tracer.Start(ctx, "multiply")
	defer span.End()

	time.Sleep(125 * time.Millisecond)
	log.Printf("multiplication is fun!")

	return &pb.CalculationResult{Result: r.GetFirstOperand() * r.GetSecondOperand()}, nil
}

func divide(ctx context.Context, r *pb.BinaryOperation) (*pb.CalculationResult, error) {
	tracer := otel.GetTracerProvider().Tracer("tracer")
	_, span := tracer.Start(ctx, "divide")
	defer span.End()

	time.Sleep(250 * time.Millisecond)
	log.Printf("division is tough")

	time.Sleep(5 * time.Second)

	return &pb.CalculationResult{Result: r.GetFirstOperand() / r.GetSecondOperand()}, nil	
}

func main() {
	tracerProvider, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		} 
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	grpcEndpoint := fmt.Sprintf(":%s", port)
	log.Printf("gRPC endpoint [%s]", grpcEndpoint)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)
	pb.RegisterCalculatorServer(grpcServer, NewServer())

	listen, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Starting: gRPC Listener [%s]\n", grpcEndpoint)
	log.Fatal(grpcServer.Serve(listen))
}
