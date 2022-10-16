package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/senpainikolay/ClientService/client"
	"github.com/senpainikolay/ClientService/structs"
)

func main() {
	rand.Seed(time.Now().UnixMilli())
	c := client.Client{-1, structs.MenuGet{}}
	clientIdCounter := client.ClientIdCounter{0, sync.Mutex{}}

	c.Work(&clientIdCounter)

}
