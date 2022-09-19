# Stilt assessment

### First-in-first-out
In this strategy, I append the resources to the relative slices and read them
from the start, this behavior is kind of similar to queue.

Below are some global variables that I have used to get it done
1. `freeCouriers []*Courier` in the start all the couriers are free.
2. `arrivedCouriers []*Courier` after arrival the courier falls in arrivedCouriers.
3. `readyOrders []*Order` prepared orders adds into the readyOrders.
4. `courierTotalWaitTime` after every delivery the courier total wait time gets updated
5. `orderTotalWaitTime` after picking up the order total wait time gets updated


The code works as follows:

There are two objects/types, one Order and another Courier.

- The code first reads the JSON files to collect the data and store them
    in the global slices (these can be created locally and passed in
    channels or params as well).
- After that code iterates over the orders and sleeps after receiving every 2 orders.
    In the meantime, the code passes the WaitGroup param to go routine functions as
    well to avoid the calculation of statistics at the end while there are orders still
    processing in the background.
- First, the process call takes place which calls the orderReceived
    orderReceived calls courierDispatched in a go-routine, because we don't want to wait
    till the order is ready.
    Then the code sleeps for the number of seconds that it takes to prepare
    Then the orderPrepared call takes place which adds the order into the readyOrders
    And once an order is prepared it is added to the readyOrders and the orderPickup
    call takes place
- In the courierDispatched call the function dispatches the first courier of the slice
    and sleeps for the number of seconds it takes to arrive.
    Once a courier is arrived the courierArrived call happens.
    From there the orderPickup call happens where an arrived courier gets an order from
    the ready ones.
- So the orderPickup is called from two functions anytime an order is ready or the courier 
    has arrived. In both cases, we want to check if there are any orders available in the
    ready orders and if any courier is available in the arrived couriers. consider an order
    delivered and put the courier back into the freeCouriers slice.
- Once we are done with all the orders we calculate the statistics from the stored data.

I have used Mutexes in the code to lock the resources while reading or writing
to avoid race conditions. I have used WaitGroup to wait for the go routine functions to end before calculating the stats
So there are no orders processing in the background.


### Matched
In this strategy, there are a few points to be careful about. One is dispatching 
the courier since we only want to dispatch the matched one. The second point is while picking up the orders, the courier should pick up the matching order only.
The match case can be assumed in various ways like matching by category matching by ID every courier has a list of orders to deliver etc.
But I've assumed that the courier arrival time should match the order prepTime. Just to make the simulation run a little longer.

The code works as follows:

There are two objects/types, one Order and another Courier.

- The code first reads the JSON files to collect the data and store them
    in the global slices (these can be created locally and passed in
    channels or params as well).
- After that code iterates over the orders and sleeps after receiving every 2 orders.
    In the meantime, the code passes the WaitGroup param to go routine functions as
    well to avoid the calculation of statistics at the end while there are orders still
    processing in the background.
- First the process call takes place which calls the orderReceived
    orderReceived calls courierDispatched in a go-routine, because we don't want to wait
    till the order is ready.
    Then the code sleeps for the number of seconds that it takes to prepare
    Then the orderPrepared call takes place which adds the order into the readyOrders
    And once an order is prepared it is added to the readyOrders and the orderPickup
    call takes place.
- In the courierDispatched call, the function calls the find function to find the matched courier 
    and dispatch it. If there isn't one then it' probably couriers are busy. The code sleeps for
    the number of seconds it takes to arrive. Once a courier is arrived the courierArrived call happens.
    From there the orderPickup call happens where an arrived courier gets an order from
    the ready ones.
- So the orderPickup is called from two functions anytime an order is ready or the courier 
    has arrived. In both cases, we want to check if there are matched orders available in the
    ready orders and if any courier is available in the arrived couriers. An order is Consider
    delivered and put the courier back into the freeCouriers slice.
- Once we are done with all the orders we calculate the statistics from the stored data.

I have used Mutexes in the code to lock the resources while reading or writing
to avoid race conditions.

