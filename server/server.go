package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/crcsmnky/grpc-calculator/proto"
	"google.golang.org/grpc"

	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	serviceName = "calculator"
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

	span, _ := tracer.StartSpanFromContext(ctx, "Calculate")
	defer span.Finish()

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
	span, _ := tracer.StartSpanFromContext(ctx, r.GetOperation().String())
	defer span.Finish()

	return &pb.CalculationResult{Result: r.GetFirstOperand() + r.GetSecondOperand()}, nil
}

func subtract(ctx context.Context, r *pb.BinaryOperation) (*pb.CalculationResult, error) {
	span, _ := tracer.StartSpanFromContext(ctx, r.GetOperation().String())
	defer span.Finish()

	return &pb.CalculationResult{Result: r.GetFirstOperand() - r.GetSecondOperand()}, nil
}

func multiply(ctx context.Context, r *pb.BinaryOperation) (*pb.CalculationResult, error) {
	span, _ := tracer.StartSpanFromContext(ctx, r.GetOperation().String())
	defer span.Finish()

	time.Sleep(125 * time.Millisecond)

	return &pb.CalculationResult{Result: r.GetFirstOperand() * r.GetSecondOperand()}, nil
}

func divide(ctx context.Context, r *pb.BinaryOperation) (*pb.CalculationResult, error) {
	span, _ := tracer.StartSpanFromContext(ctx, r.GetOperation().String())
	defer span.Finish()

	time.Sleep(250 * time.Millisecond)

	return &pb.CalculationResult{Result: r.GetFirstOperand() + r.GetSecondOperand()}, nil
}

func main() {
	tracer.Start(tracer.WithDebugMode(true))
	defer tracer.Stop()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	grpcEndpoint := fmt.Sprintf(":%s", port)
	log.Printf("gRPC endpoint [%s]", grpcEndpoint)

	ui := grpctrace.UnaryServerInterceptor(
		grpctrace.WithServiceName(serviceName),
	)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(ui))
	pb.RegisterCalculatorServer(grpcServer, NewServer())

	listen, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Starting: gRPC Listener [%s]\n", grpcEndpoint)
	log.Fatal(grpcServer.Serve(listen))
}
