package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/senpainikolay/ClientService/structs"
)

type Client struct {
	ClientId int
	ResInfo  structs.MenuGet
}

type ClientIdCounter struct {
	IdCounter int
	Mutex     sync.Mutex
}

func (c *Client) Work(cic *ClientIdCounter, address string) {
	c.GetId(cic)
	c.RequestMenu(address)
	c.GenerateOrdersAndSendToOM(address)
	time.Sleep(time.Duration(rand.Intn(40)+70) * 100 * time.Millisecond)
	go func() { c.Work(cic, address) }()
	time.Sleep(10 * time.Millisecond)
}

func (c *Client) GenerateOrdersAndSendToOM(address string) {
	resIdSlice := c.GenerateRandomRestaurantIds()
	log.Println(resIdSlice)

	var orders structs.Orders
	orders.ClientId = c.ClientId
	var wg sync.WaitGroup
	wg.Add(len(resIdSlice))
	for _, id := range resIdSlice {
		idx_id := id
		go func() {
			ord := c.GenerateOneOrder(idx_id)
			ord.RestaurantId = c.ResInfo.RestaurantsData[idx_id].RestaurantId
			orders.Orders = append(orders.Orders, ord)
			wg.Done()
		}()
	}
	wg.Wait()
	log.Println(orders)
	SendOrderToOM(&orders, address)

}

func SendOrderToOM(ords *structs.Orders, address string) {

	postBody, _ := json.Marshal(ords)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://"+address+"/order", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var clientRes structs.ClientResponse
	if err := json.Unmarshal([]byte(body), &clientRes); err != nil {
		panic(err)
	}

	log.Println(clientRes)

}

func (c *Client) GetId(cic *ClientIdCounter) {
	cic.Mutex.Lock()
	cic.IdCounter += 1
	c.ClientId = cic.IdCounter
	cic.Mutex.Unlock()

}

func (c *Client) GenerateRandomRestaurantIds() []int {
	var generatedResIds []int
	for k := 0; k < c.ResInfo.Restaurants; k++ {
		rn := rand.Intn(c.ResInfo.Restaurants)
		rnBool := false
		for j := 0; j < len(generatedResIds); j++ {
			if rn == generatedResIds[j] {
				rnBool = true
				break
			}
		}
		if rnBool {
			continue
		}

		generatedResIds = append(generatedResIds, rn)

	}
	return generatedResIds

}

func (c *Client) GenerateOneOrder(idx_resId int) structs.Order {
	items := make([]int, rand.Intn(8)+1)
	maxWaitInt := 0
	for i := range items {
		items[i] = rand.Intn(c.ResInfo.RestaurantsData[idx_resId].MenuItems) + 1
		if c.ResInfo.RestaurantsData[idx_resId].Menu[items[i]-1].PreparationTime > maxWaitInt {
			maxWaitInt = c.ResInfo.RestaurantsData[idx_resId].Menu[items[i]-1].PreparationTime
		}
	}
	priority := rand.Intn(5) + 1

	return structs.Order{
		RestaurantId: -1,
		Items:        items,
		Priority:     priority,
		MaxWait:      float64(maxWaitInt) * 1.8,
		CreatedTime:  time.Now().UnixMilli(),
	}

}

func (c *Client) RequestMenu(address string) {
	resp, err := http.Get("http://" + address + "/menu")
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var menuInfo structs.MenuGet
	if err := json.Unmarshal([]byte(body), &menuInfo); err != nil {
		panic(err)
	}
	c.ResInfo = menuInfo
	// log.Println(menuInfo)

}
