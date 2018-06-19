package main

import (
	"flag"
	"log"
	"io/ioutil"
	"os"
	"fmt"

	pb "../api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	backend := flag.String("b", "localhost:8080", "address of backend")
	output := flag.String("o", "output.wav", "output file name")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("usage:\n\t%s \"text to speech\"\n", os.Args[0])
		os.Exit(1)
	}

	conn, err := grpc.Dial(*backend, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to %s: %v", *backend, err)
	}

	defer conn.Close()

	client := pb.NewTextToSpeechClient(conn)
	text := &pb.Text{ Text: flag.Arg(0)}
	res, err := client.Say(context.Background(), text)
	if err != nil {
		log.Fatalf("Couldn't say %s: %v", text.Text, err)
	}

	if err := ioutil.WriteFile(*output, res.Speech, 0666); err != nil {
		log.Fatalf("Couldn't write output file %s: %v", *output, err)
	}
}