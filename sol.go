package main

import (
	"encoding/json"
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

// An order for kitchen holds relative information
type Order struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	PrepTime int64  `json:"prepTime"`
	State    string
	ReadyAt  time.Time
	PickedAt time.Time
}

// An courier for delivery holds relative information
type Courier struct {
	Name       string `json:"name"`
	ArriveTime int64  `json:"arrivalTime"`
	ArrivedAt  time.Time
	PickedAt   time.Time
}

// Acknowledges the receive of order and calls orderPrepared
// also calls for the courier dispatch
func (o *Order) orderReceived(mx *sync.Mutex, wg *sync.WaitGroup) {
	// log order received
	log.Info("Order received: ", o.Id)
	// call the courierDispatched
	wg.Add(1)
	go courierDispatched(mx, wg)
	// sleep till order is ready
	time.Sleep(time.Duration(o.PrepTime) * time.Second)
	// call orderPrepared
	o.orderPrepared(mx)
}

// Acknowledges the prepare of order and calls for pickup
// updates the readyAt field of the order
func (o *Order) orderPrepared(mx *sync.Mutex) {
	o.ReadyAt = time.Now()
	// log order prepared
	log.Info("Order prepared: ", o.Id)
	// add the order into prepared orders
	mx.Lock()
	readyOrders = append(readyOrders, o)
	mx.Unlock()
	orderPickedUp(mx)

}

// courierDispatched dispatches the 1st courier of the freeCouriers
// and calls for the courier to be arrived
// after the sleep. Then calls for courier arrived
func courierDispatched(mx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	// move a courier from free to busy
	if len(freeCouriers) < 1 {
		return
	}

	mx.Lock()
	courier := freeCouriers[0]
	freeCouriers = append(freeCouriers[:0], freeCouriers[1:]...)
	mx.Unlock()

	// log courier dispatched
	log.Info("Courier dispatched: ", courier.Name)

	// wait till courier arrives
	time.Sleep(time.Duration(courier.ArriveTime) * time.Second)

	// call CourierArrived
	courier.courierArrived(mx)
}

// Acknowledges the arrival of courier and calls for pickup
func (c *Courier) courierArrived(mx *sync.Mutex) {
	// log courier arrived
	log.Info("Courier arrived: ", c.Name)

	c.ArrivedAt = time.Now()
	mx.Lock()
	arrivedCouriers = append(arrivedCouriers, c)
	mx.Unlock()

	// call orderPickedUp
	orderPickedUp(mx)
}

// Exits if there is no order to pickup or no arrivedCourier
// it removes the courier from arrivedCourier after picking and adds
// him back in the free couriers. It also removes the order from Orders
func orderPickedUp(mx *sync.Mutex) {
	// defer wg.Done()
	// pick an order from ready orders
	if len(arrivedCouriers) < 1 || len(readyOrders) < 1 {
		return
	}

	mx.Lock()
	// remove the courier from arrivedCouriers
	courier := arrivedCouriers[0]
	arrivedCouriers = arrivedCouriers[1:]

	// remove the order from readyOrders
	order := readyOrders[0]
	readyOrders = readyOrders[1:]

	orderTotalWaitTime += time.Since(order.ReadyAt)
	courierTotalWaitTime += time.Since(courier.ArrivedAt)

	// log order picked up
	log.Infof("%s Order Picked by %s courier", order.Id, courier.Name)

	// add the courier back to the freeCouriers after the pickup
	freeCouriers = append(freeCouriers, courier)
	mx.Unlock()
}

// Starts the order processing flow
func (o *Order) process(mx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	// order received
	o.orderReceived(mx, wg)
}

func main() {
	// read the orders JSON into orders slice
	rawOrders, err := ioutil.ReadFile("./dispatch_orders.json")
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
	rawCouriers, err := ioutil.ReadFile("./couriers.json")
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
