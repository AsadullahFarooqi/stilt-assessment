package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("KitchenOrdersDelivery")

var freeCouriers []*Courier
var arrivedCouriers []*Courier

var readyOrders []*Order

var courierTotalWaitTime time.Duration
var orderTotalWaitTime time.Duration

type Order struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	PrepTime int64  `json:"prepTime"`
	State    string
	ReadyAt  time.Time
	PickedAt time.Time
}

type Courier struct {
	Name       string `json:"name"`
	ArriveTime int64  `json:"arrivalTime"`
	ArrivedAt  time.Time
	PickedAt   time.Time
}

func findCourier(couriers []*Courier, prepTime int64) (int, error) {
	for i, c := range couriers {
		if c.ArriveTime == prepTime {
			return i, nil
		}
	}
	return -1, fmt.Errorf("courier not found")
}

func findOrder(orders []*Order, prepTime int64) (int, error) {
	for i, c := range orders {
		if c.PrepTime == prepTime {

			return i, nil
		}
	}
	return -1, fmt.Errorf("order not found")
}

func (o *Order) orderReceived(mx *sync.Mutex, wg *sync.WaitGroup) {
	// log order received
	log.Info("Order received: ", o.Id)
	// call the courierDispatched
	wg.Add(1)
	go courierDispatched(mx, wg, o)
	// sleep till order is ready
	time.Sleep(time.Duration(o.PrepTime) * time.Second)
	// call orderPrepared
	o.orderPrepared(mx)
}

func (o *Order) orderPrepared(mx *sync.Mutex) {
	o.ReadyAt = time.Now()
	// log order prepared
	log.Info("Order prepared: ", o.Id)
	// add the order into prepared orders
	mx.Lock()
	readyOrders = append(readyOrders, o)
	mx.Unlock()

	orderPickedUp(mx, o.PrepTime)

}

func courierDispatched(mx *sync.Mutex, wg *sync.WaitGroup, o *Order) {
	defer wg.Done()
	// move a courier from free to busy
	if len(freeCouriers) < 1 {
		return
	}

	mx.Lock()
	courierIndex, err := findCourier(freeCouriers, o.PrepTime)
	if err != nil {
		log.Warning("courier not found they must be busy, Error: %s", err.Error())
		mx.Unlock()
		return
	}

	courier := freeCouriers[courierIndex]

	if courierIndex+1 == len(freeCouriers) {
		freeCouriers = freeCouriers[:courierIndex]
	} else {
		freeCouriers = append(freeCouriers[:courierIndex], freeCouriers[courierIndex+1:]...)
	}
	mx.Unlock()

	// log courier dispatched
	log.Info("Courier dispatched: ", courier.Name)

	// wait till courier arrives
	time.Sleep(time.Duration(courier.ArriveTime) * time.Second)

	// call CourierArrived
	courier.courierArrived(mx)
}

func (c *Courier) courierArrived(mx *sync.Mutex) {
	// log courier arrived
	log.Info("Courier arrived: ", c.Name)

	c.ArrivedAt = time.Now()
	mx.Lock()
	arrivedCouriers = append(arrivedCouriers, c)
	mx.Unlock()

	// call orderPickedUp
	orderPickedUp(mx, c.ArriveTime)
}

func orderPickedUp(mx *sync.Mutex, matchingPoint int64) {
	// pick an order from ready orders
	if len(arrivedCouriers) < 1 || len(readyOrders) < 1 {
		return
	}

	mx.Lock()
	// find a courier for the delivery
	courierIndex, err := findCourier(arrivedCouriers, matchingPoint)
	if err != nil {
		log.Warning("couriers must be on the way, Error: ", err.Error())
		mx.Unlock()
		return
	}

	// take the 1st available matching order for delivery
	orderIndex, err := findOrder(readyOrders, matchingPoint)
	if err != nil {
		log.Warning("order must preparing, Error: ", err.Error())
		mx.Unlock()
		return
	}

	// remove the courier from the arrivedCouriers
	courier := arrivedCouriers[courierIndex]
	rmvFromArrivedCouriers(courierIndex)

	order := readyOrders[orderIndex]
	// remove the order from readyOrders
	rmvFromReadyOrders(orderIndex)

	orderTotalWaitTime += time.Since(order.ReadyAt)
	courierTotalWaitTime += time.Since(courier.ArrivedAt)

	// log order picked up
	log.Infof("%s Order Picked by %s courier", order.Id, courier.Name)

	// add the courier back to the freeCouriers after the pickup
	freeCouriers = append(freeCouriers, courier)
	mx.Unlock()
}

func (o *Order) process(mx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	// order received
	o.orderReceived(mx, wg)
}

func rmvFromArrivedCouriers(index int) {
	if int(index+1) == len(arrivedCouriers) {
		arrivedCouriers = arrivedCouriers[:index]
	} else {
		arrivedCouriers = append(arrivedCouriers[:index], arrivedCouriers[index+1:]...)
	}
}

func rmvFromReadyOrders(index int) {
	if index+1 == len(readyOrders) {
		readyOrders = readyOrders[:index]
	} else {
		readyOrders = append(
			readyOrders[:index],
			readyOrders[index+1:]...)
	}
}

func main() {
	// read the orders JSON into orders slice
	rawOrders, err := ioutil.ReadFile("../dispatch_orders copy.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var orders []*Order

	// var orders []Order
	err = json.Unmarshal(rawOrders, &orders)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	// read the couriers JSON into the couriers slice
	rawCouriers, err := ioutil.ReadFile("../couriers.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// var couriers []Courier
	err = json.Unmarshal(rawCouriers, &freeCouriers)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	numOfOrders := len(orders)

	var m sync.Mutex
	var wg sync.WaitGroup

	// iterate over orders
	for i, order := range orders {
		wg.Add(1)
		go order.process(&m, &wg)

		// sleep after every 2 orders into the on-going-orders slice
		if i != 0 && i%2 == 0 {
			time.Sleep(1 * time.Second)
		}
	}

	wg.Wait()

	log.Info("Average food wait time: ", float64(orderTotalWaitTime.Milliseconds())/float64(numOfOrders))
	log.Info("Average courier wait time: ", float64(courierTotalWaitTime.Milliseconds())/float64(len(freeCouriers)))

}
