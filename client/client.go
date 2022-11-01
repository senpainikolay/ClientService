package client

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	ordersStatus := c.GenerateOrdersAndSendToOM(address, cic)
	ordersLen := len(ordersStatus.Orders)
	if ordersLen != 0 {
		log.Printf("Client ID %v order send to FoodOrdering", c.ClientId)
		var wg sync.WaitGroup
		wg.Add(ordersLen)
		for _, or := range ordersStatus.Orders {
			order := or
			go func() {
				timeToSleep := time.Duration(order.EstimatedWaitingTime)
				for {
					time.Sleep(timeToSleep * 100 * time.Millisecond)
					respOrdStatus := SendGetOrderStatusToRes(order.RestaurantAddress, order.OrderId)
					if respOrdStatus.IsReady {
						// Logic For Order Ready
						break
					} else {
						timeToSleep = time.Duration(respOrdStatus.EstimatedWaitingTime)
					}
				}
				wg.Done()
			}()
		}

		wg.Wait()
		log.Printf("Client ID %v Succesfully recieved the order back", c.ClientId)
	} else {
		// Restaurant full, cancels the Client.
		// log.Println("RIPP")
		time.Sleep(100 * 50 * time.Millisecond)
	}
	// Generate new Client
	go func() { c.Work(cic, address) }()

}

func SendGetOrderStatusToRes(address string, id int) *structs.ClientOrderStatus {

	resp, err := http.Get("http://" + address + "/v2/order/" + fmt.Sprint(id))
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var ordInfo structs.ClientOrderStatus
	if err := json.Unmarshal([]byte(body), &ordInfo); err != nil {
		panic(err)
	}
	return &ordInfo

}

func (c *Client) GenerateOrdersAndSendToOM(address string, cic *ClientIdCounter) *structs.ClientResponse {
	resIdSlice := c.GenerateRandomRestaurantIds()
	// log.Println(resIdSlice)
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
	return SendOrderToOM(&orders, address, cic)

}

func SendOrderToOM(ords *structs.Orders, address string, cic *ClientIdCounter) *structs.ClientResponse {

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

	if len(clientRes.Orders) == 0 {
		decreaseId(cic)
	}

	return &clientRes

}

func (c *Client) GetId(cic *ClientIdCounter) {
	cic.Mutex.Lock()
	cic.IdCounter += 1
	c.ClientId = cic.IdCounter
	cic.Mutex.Unlock()

}

func decreaseId(cic *ClientIdCounter) {
	cic.Mutex.Lock()
	cic.IdCounter -= 1
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
