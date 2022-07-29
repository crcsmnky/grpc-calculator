package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

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
	log.Println("[server:Calculate] Started")
	if ctx.Err() == context.Canceled {
		return &pb.CalculationResult{}, fmt.Errorf("client cancelled: abandoning")
	}

	switch r.GetOperation() {			
		case pb.Operation_ADD:
			return &pb.CalculationResult{
				Result: r.GetFirstOperand() + r.GetSecondOperand(),
			}, nil
		case pb.Operation_SUBTRACT:
			return &pb.CalculationResult{
				Result: r.GetFirstOperand() - r.GetSecondOperand(),
			}, nil
		case pb.Operation_MULTIPLY:
			return &pb.CalculationResult{
				Result: r.GetFirstOperand() * r.GetSecondOperand(),
			}, nil
		case pb.Operation_DIVIDE:
			return &pb.CalculationResult{
				Result: r.GetFirstOperand() / r.GetSecondOperand(),
			}, nil			
		default:
			return &pb.CalculationResult{}, fmt.Errorf("undefined operation")
	}
}

func main() {
	tracer.Start()
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
