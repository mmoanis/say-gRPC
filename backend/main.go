package main

import (
	"flag"
	"fmt"
	"net"
	"os/exec"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	pb "../api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func main() {
	port := flag.Int("p", 8080, "port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logrus.Fatalf("Couldn't listen to port %d: %v", *port, err)
	}

	logrus.Infof("Listening to %d", *port)

	s := grpc.NewServer()
	pb.RegisterTextToSpeechServer(s, Server{})

	err = s.Serve(lis)
	if err != nil {
		logrus.Fatalf("Couldn't serve: %v", err)
	}
}


type Server struct {
}

func (Server) Say(ctx context.Context, text *pb.Text) (*pb.Speech, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get peer from ctx")
	}
	if pr.Addr == net.Addr(nil) {
		return nil, fmt.Errorf("failed to get peer address")
	}

	logrus.Infof("Serving RPC from %s", pr.Addr)

	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("Couldn't create temp file")
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("Couldn't close file %s: %v", f.Name(), err)
	}
	
	cmd := exec.Command("flite", "-t", text.Text, "-o", f.Name())
	if data, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("flite failed: %s - %v", data, err)
	}

	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("Couldn't read tmp file: %v", err)
	}

	return &pb.Speech{Speech: data}, nil
}
