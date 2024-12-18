
// Server 3: gRPC server (s3)
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "protobuf"
)

type s3Server struct {
	pb.UnimplementedS3ServiceServer
}

func (s *s3Server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	result := req.Value1 + req.Value2
	return &pb.AddResponse{Result: result}, nil
}

func main() {
	
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterS3ServiceServer(grpcServer, &s3Server{})

	log.Println("Starting gRPC server on :50052")
	log.Fatal(grpcServer.Serve(listener))
}
