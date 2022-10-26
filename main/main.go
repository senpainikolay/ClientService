package main

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/senpainikolay/ClientService/client"
	"github.com/senpainikolay/ClientService/structs"
)

func main() {
	r := mux.NewRouter()
	rand.Seed(time.Now().UnixMilli())
	conf := structs.GetConf()
	var clients []client.Client
	clientIdCounter := client.ClientIdCounter{0, sync.Mutex{}}
	for i := 0; i < 2; i++ {
		clients = append(clients, client.Client{-1, structs.MenuGet{}})
	}
	for i := 0; i < 2; i++ {
		id := i
		go func() {

			clients[id].Work(&clientIdCounter, conf.OMAddress)

		}()

	}
	http.ListenAndServe(":"+conf.Port, r)
}
