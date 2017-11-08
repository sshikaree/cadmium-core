package main

import (
	"log"

	"github.com/sshikaree/cadmium-core/common/database"
	"github.com/sshikaree/cadmium-core/common/rpc"
	"github.com/sshikaree/cadmium-core/common/webinterface"
)

func main() {
	var err error
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	database.DB, err = database.NewSqliteProvider("./db/cadmium.db")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		webinterface.StartWebInterface("7000")
	}()

	// rpc.StartLocalSocketRPCServer()
	rpc.ConnectToActiveAccounts()
	rpc.StartWebsocketRPCServer("7001")

}
