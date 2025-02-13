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

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"time"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"github.com/google/uuid"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "world"
)


var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	status = int32(1)
)

/*func convertStatus(s int) proto.UpdateNodeAttestStatusRequest_Code {
	m := map[int]proto.UpdateNodeAttestStatusRequest_Code{
		0: proto.UpdateNodAttestStatusRequest_ATTEST_STATUS_SUCEESS,
		1: proto.UpdateNodeAttestStatusRequest_ATTEST_STATUS_FAIL,
	}
	return m[s]
}
*/
func main() {
	flag.Parse()

	caCert, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		log.Fatalf("failed to read CA certificate: %s", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create the credentials and return it
	config := &tls.Config{
        RootCAs:      caCertPool,
    }
	creds := credentials.NewTLS(config)

	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNodeAttestationManagerServiceClient(conn)

	//Generate a random UUID
	reqUuid := uuid.New().String()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UpdateNodeAttestationStatus(ctx, &pb.UpdateNodeAttestStatusRequest{Code: status, systemuuid:reqUuid})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
