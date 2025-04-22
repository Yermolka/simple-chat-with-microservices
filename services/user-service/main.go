package main

import (
	"log"
	"net"
	"user-service/internal/grpc"
	"user-service/internal/http"
	"user-service/internal/repository"
)

func main() {
	repo := repository.NewRepository()
	server := http.NewServer(repo)
	go server.ListenAndServe()

	grpcServer := grpc.NewServer(repo)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer.Serve(lis)

}
