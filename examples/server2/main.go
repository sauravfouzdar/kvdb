// Server 2: gRPC server (s2)
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "path/to/protobuf"
)

type s2Server struct {
	pb.UnimplementedS2ServiceServer
}

func (s *s2Server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	// Forward request to s3
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	s3Client := pb.NewS3ServiceClient(conn)
	s3Req := &pb.AddRequest{
		Value1: req.Value1,
		Value2: req.Value2,
	}

	s3Res, err := s3Client.Add(context.Background(), s3Req)
	if err != nil {
		return nil, err
	}

	return &pb.AddResponse{Result: s3Res.Result}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterS2ServiceServer(grpcServer, &s2Server{})

	log.Println("Starting gRPC server on :50051")
	log.Fatal(grpcServer.Serve(listener))
}
