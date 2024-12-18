// Server 1: REST API server (s1)
package main

import (
	//"bytes"
	"encoding/json"
	//"fmt"
	// "io/ioutil"
	"context"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Import the protobuf package for gRPC
	pb "github.com/sauravfouzdar/kvdb/examples/proto"
)

type AddRequest struct {
	Value1 int `json:"value1"`
	Value2 int `json:"value2"`
}

type AddResponse struct {
	Result int `json:"result"`
}

func main() {
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var req AddRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Failed to decode request", http.StatusBadRequest)
			return
		}

		// Call s2 (gRPC server)
		conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			http.Error(w, "Failed to connect to gRPC server", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		client := pb.NewS2ServiceClient(conn)
		grpcReq := &pb.AddRequest{
			Value1: int32(req.Value1),
			Value2: int32(req.Value2),
		}

		grpcRes, err := client.Add(context.Background(), grpcReq)
		if err != nil {
			http.Error(w, "gRPC call failed", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(AddResponse{Result: int(grpcRes.Result)})
	})

	log.Println("Starting REST API server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
