package main

import (
	"context"
	"fmt"
	"io"
	"log"
	api "randsig/pkg/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	connection, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	client := api.NewRandomSignalerClient(connection)
	stream, err := client.RandSignal(context.Background(), &api.RandSignalRequest{})
	if err != nil {
		log.Fatal(err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(message.SessionId, message.Frequency, message.CurrentTimestamp)
	}
}
