package client

import (
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

func (c *Client) Work(cic *ClientIdCounter) {
	c.GetId(cic)
	c.RequestMenu()
	c.GenerateOrders()

}

func (c *Client) GenerateOrders() {
	resIdSlice := c.GenerateRandomRestaurantIds()
	var orders structs.Orders
	orders.ClientId = c.ClientId
	var wg sync.WaitGroup
	wg.Add(len(resIdSlice))
	for i := range resIdSlice {
		id := i
		go func() {
			ord := c.GenerateOneOrder(id)
			ord.RestaurantId = id + 1
			orders.Orders = append(orders.Orders, ord)
			wg.Done()
		}()
	}
	wg.Wait()
	log.Println(orders)

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

func (c *Client) GenerateOneOrder(resId int) structs.Order {
	items := make([]int, rand.Intn(10)+1)
	maxWaitInt := 0
	for i := range items {
		items[i] = rand.Intn(c.ResInfo.RestaurantsData[resId].MenuItems) + 1
		if c.ResInfo.RestaurantsData[resId].Menu[items[i]-1].PreparationTime > maxWaitInt {
			maxWaitInt = c.ResInfo.RestaurantsData[resId].Menu[items[i]-1].PreparationTime
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

func (c *Client) RequestMenu() {
	resp, err := http.Get("http://localhost:5000/menu")
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
