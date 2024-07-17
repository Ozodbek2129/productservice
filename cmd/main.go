package main

import (
	"log"
	"net"
	"product/config"
	"product/service"
	"product/storage/postgres"
	"product/genproto/ProductService"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp",config.Load().PRODUCT_SERVICE)
	if err!=nil{
		log.Fatal(err)
	}
	defer listener.Close()

	db,err:=postgres.ConnectDB()
	if err!=nil{
		log.Fatal(err)
	}
	defer db.Close()

	productservice:=service.NewArtisanConnectService(db)
	server:=grpc.NewServer()

	ProductService.RegisterProductServiceServer(server,productservice)

	log.Printf("Server is listening on port %s\n",config.Load().PRODUCT_SERVICE)
	if err=server.Serve(listener); err!=nil{
		log.Fatal(err)
	}
}
