package main

import (
	"log"
	"net"
	msg"github.com/Prameesh-P/SHOPRIX/messages"

	"google.golang.org/grpc"
)


const(

	port = ":6000"
)

type Server struct{
	msg.AuthentifiationServiceServer
}

func main() {
	LoadEnv()
	ConnectToDb()
	SyncDB()
	lis, err := net.Listen("tcp",port)
	if err!=nil{
		log.Fatalf("Error got while listening on port:%s error is >>> %v",port,err)
	}

	grpcServer:=grpc.NewServer()
	msg.RegisterAuthentifiationServiceServer(grpcServer,&Server{})
	log.Printf("Server is starting on the port %v",port)
	if err:=grpcServer.Serve(lis);err!=nil{
			log.Fatalf("Cannot start grpc server got error %v",err)
	}	
}