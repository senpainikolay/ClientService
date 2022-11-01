package structs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Food struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	PreparationTime  int    `json:"preparation_time"`
	Complexity       int    `json:"complexity"`
	CookingApparatus string `json:"cooking_apparatus"`
}

type Menu []Food

type RestaurantData struct {
	RestaurantId int     `json:"restaurant_id"`
	Name         string  `json:"name"`
	MenuItems    int     `json:"menu_items"`
	Menu         Menu    `json:"menu"`
	Rating       float64 `json:"rating"`
}

type MenuGet struct {
	Restaurants     int              `json:"restaurants"`
	RestaurantsData []RestaurantData `json:"restaurants_data"`
}

type Order struct {
	RestaurantId int     `json:"restaurant_id"`
	Items        []int   `json:"items"`
	Priority     int     `json:"priority"`
	MaxWait      float64 `json:"max_wait"`
	CreatedTime  int64   `json:"created_time"`
}

type Orders struct {
	ClientId int     `json:"client_id"`
	Orders   []Order `json:"orders"`
}

type OMResponse struct {
	RestaurantId         int     `json:"restaurant_id"`
	RestaurantAddress    string  `json:"restaurant_address"`
	OrderId              int     `json:"order_id"`
	EstimatedWaitingTime float64 `json:"estimated_waiting_time"`
	CreatedTime          int64   `json:"created_time"`
	RegisteredTime       int64   `json:"registered_time"`
}

type ClientResponse struct {
	OrderId int          `json:"order_id"`
	Orders  []OMResponse `json:"orders"`
}

type Conf struct {
	Port      string `json:"port"`
	OMAddress string `json:"om_address"`
}

type ClientOrderStatus struct {
	OrderId              int              `json:"order_id"`
	IsReady              bool             `json:"is_ready"`
	EstimatedWaitingTime float64          `json:"estimated_waiting_time"`
	Priority             int              `json:"priority"`
	MaxWait              float64          `json:"max_wait"`
	CreatedTime          int64            `json:"created_time"`
	RegisteredTime       int              `json:"registered_time"`
	PreparedTime         int64            `json:"prepared_time"`
	CookingTime          int64            `json:"cooking_time"`
	CookingDetails       []CookingDetails `json:"cooking_details"`
}

type CookingDetails struct {
	CookId int `json:"cook_id"`
	FoodId int `json:"food_id"`
}

func GetConf() *Conf {
	jsonFile, err := os.Open("configurations/Conf.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var conf Conf
	json.Unmarshal(byteValue, &conf)
	return &conf

}
