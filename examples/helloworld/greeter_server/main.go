/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"net"
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedNodeAttestationManagerServiceServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) UpdateNodeAttestationStatus (_ context.Context, in *pb.UpdateNodeAttestStatusRequest) (*pb.UpdateNodeAttestStatusResponse, error) {
	log.Printf("Received: %v", in.Code)
	return &pb.UpdateNodeAttestStatusResponse{Message: "Hello "}, nil
}

func main() {
	flag.Parse()
	
	cert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		log.Fatalf("failed to load server key pair: %s", err)
	}

	// Create the credentials
	config := &tls.Config{
	Certificates: []tls.Certificate{cert},
	ClientAuth:   tls.NoClientCert,
    }

	creds := credentials.NewTLS(config)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterNodeAttestationManagerServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
