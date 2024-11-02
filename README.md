# Stilt assessment

As the assessment discouraged the use of production technologies thats why I didn't write the Docker file.

## running the code
The first parent directory has the first-in-first-out strategy code. To run that code
you can just hit `ctrl+F5` in the VS-Code. If you don't have the vs-code setup then
open a terminal change your directory to the assessment directory and run `go run sol.go` command.

To run the matched_strategy change your directory to the <assessment>/matched_strategy and run the same
command.

if you don't have the Golang setup follow [this link](https://go.dev/doc/install)

### First-in-first-out
In this strategy, I append the resources to the relative slices and read them
from the start, this behavior is kind of similar to queue.

The code works as follows:

- The code first reads the JSON files to collect the data and store them
    in the global slices (these can be created locally and passed in
    channels or params as well).
- After that code iterates over the orders and sleeps after receiving every 2 orders.
    In the meantime, the code passes the WaitGroup param to go routine functions as
    well to avoid the calculation of statistics at the end while there are orders still
    processing in the background.
- First, the process call takes place which calls the orderReceived
    orderReceived calls courierDispatched in a go-routine because we don't want to wait
    till the order is ready.
    Then the code sleeps for the number of seconds that it takes to prepare
    Then the orderPrepared call takes place which adds the order to the readyOrders
    And once an order is prepared it is added to the readyOrders and the orderPickup
    call takes place
- In the courierDispatched call the function dispatches the first courier of the slice
    and sleeps for the number of seconds it takes to arrive.
    Once a courier has arrived the courierArrived call happens.
    From there, the orderPickup call happens where an arrived courier gets an order from
    the ready ones.
- So the orderPickup is called from two functions anytime an order is ready or the courier 
    has arrived. In both cases, we want to check if there are any orders available in the
    ready orders and if any courier is available in the arrived couriers. consider an order
    delivered and put the courier back into the freeCouriers slice.
- Once we are done with all the orders we calculate the statistics from the stored data.


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
- First, the process call takes place which calls the orderReceived
    orderReceived calls courierDispatched in a go-routine, because we don't want to wait
    till the order is ready.
    Then the code sleeps for the number of seconds that it takes to prepare
    Then the orderPrepared call takes place which adds the order to the readyOrders
    And once an order is prepared it is added to the readyOrders and the orderPickup
    call takes place.
- In the courierDispatched call, the function calls the find function to find the matched courier 
    and dispatch it. If there isn't one then it's probably couriers who are busy. The code sleeps for
    the number of seconds it takes to arrive. Once a courier has arrived the courierArrived call happens.
    From there, the orderPickup call happens when an arrived courier gets an order from
    the ready ones.
- So the orderPickup is called from two functions anytime an order is ready or the courier 
    has arrived. In both cases, we want to check if there are matched orders available in the
    ready orders and if any courier is available in the arrived couriers. An order is considered
    delivered and puts the courier back into the freeCouriers slice.
- Once we are done with all the orders we calculate the statistics from the stored data.

A map could be used for storing the arrivedCouriers as well as the preparation time as
a key and a slice of couriers. Just like it's done for orders.


In both strategies, I have used Mutexes in the code to lock the resources while reading or writing
to avoid race conditions. I have used WaitGroup to wait for the go routine functions to end before
calculating the stats So there are no orders processing in the background.


## Improvements that could have been made 
It could have been done more structurally but my apologies due to my busy
days that's all I could do to stay up the night 🙂. Below are a few
improvements that can take place:

1. The tests aren't very good
2. error handling or edge cases handling, for example, it is possible there is no match for an order in the couriers
   in that case the program will halt and wait.
3. All these global variables can be passed as params as well
4. JSON File names should also come from the environment or command line and have default values in case of nil values.
5. A docker file may be helpful in different setups.

**Final note I can do this task in Python and Nodejs/JS as well.**
