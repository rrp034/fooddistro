// in file gatekeeper/gatekeeper.go
//
// The software suite consists of main.go,
// food, stats, gatekeeper,
// producer, foodstore, and sales packages.
//
// The GateKeeper accepts FoodPacks from inter-galactic transports
// and repacks them into FoodPacks suitable for Sales to
// distribute to the planets in the system.

package gatekeeper

import (
	"fmt"
	"time"

	"fooddistro/food"
	"fooddistro/foodstore" // HIGHLIGHT: circularque -> foodstore
	// HIGHLIGHT END
)

type GateKeeper struct {
	store *foodstore.FoodStore // HIGHLIGHT: was *circularque.CircularQue[food.FoodPack]
	// HIGHLIGHT END
	acceptChan   chan food.FoodPack
	retrieveChan chan retrieveRequest
	rejected     int
	startTime    time.Time
	endTime      time.Time

	// HIGHLIGHT: per-type/total sales counters for closing report
	soldCounts map[food.FoodType]int
	totalSold  int
	// HIGHLIGHT END
}

type retrieveRequest struct {
	response  chan food.FoodPack
	available chan bool
}

// NewGateKeeper creates and starts a GateKeeper with the given storage capacity.
func NewGateKeeper(capacity int) *GateKeeper {
	gk := &GateKeeper{
		store:        foodstore.NewFoodStore(capacity), // HIGHLIGHT: builds the sorted-list store
		acceptChan:   make(chan food.FoodPack, 100),
		retrieveChan: make(chan retrieveRequest, 100),
		rejected:     0,
		startTime:    time.Now(),
		endTime:      time.Now().Add(40 * time.Second), // 40 seconds = 40 hours simulation
		soldCounts:   make(map[food.FoodType]int),      // HIGHLIGHT: init the new sales counters
		// HIGHLIGHT END
	}
	go gk.run()
	return gk
}

// AcceptMessage requests storage space for an arriving FoodPack.
func (gk *GateKeeper) AcceptMessage(foodPack food.FoodPack) {
	gk.acceptChan <- foodPack
}

// RetrieveMessage requests the next FoodPack for sale from the GateKeeper.
func (gk *GateKeeper) RetrieveMessage() (food.FoodPack, bool) {
	respChan := make(chan food.FoodPack)
	availChan := make(chan bool)

	gk.retrieveChan <- retrieveRequest{response: respChan, available: availChan}

	available := <-availChan
	if available {
		return <-respChan, true
	}
	return food.FoodPack{}, false
}

// run is the GateKeeper's single serialized event loop.
func (gk *GateKeeper) run() {
	time.Sleep(500 * time.Millisecond) // Allow initialization (0.5 seconds = 0.5 hours)

	for gk.rejected < 5 && time.Now().Before(gk.endTime) {
		select {
		case newFood := <-gk.acceptChan:
			if !gk.store.IsFull() {
				gk.store.AcceptMessage(newFood)
				fmt.Printf("GateKeeper insert accepted %s %c\n", newFood.FoodType, newFood.FoodShipment)
			} else {
				gk.rejected++
				fmt.Printf("Rejected by GateKeeper:\n")
				fmt.Printf("%s %c\n", newFood.FoodType, newFood.FoodShipment)
				fmt.Printf("Rejected = %d. Sent to another distribution facility!\n\n", gk.rejected)
			}

		case req := <-gk.retrieveChan:
			if !gk.store.IsEmpty() {
				mgtDesiredType := food.RandomFoodType()
				// HIGHLIGHT: binary search replaces sequential FIFO removal
				foodItem, _, apology := gk.store.RetrieveDesired(mgtDesiredType)
				// HIGHLIGHT END

				fmt.Printf("Mgt Desired Food Type To Sell is: %s\n", mgtDesiredType)
				if apology != "" {
					fmt.Println(apology)
				}
				fmt.Printf("Actual type sold is: %s\n", foodItem.FoodType)
				fmt.Printf("Food pack removed by GateKeeper for shipment.\n\n")

				// HIGHLIGHT: tally sale for closing summary
				gk.soldCounts[foodItem.FoodType]++
				gk.totalSold++
				// HIGHLIGHT END

				req.available <- true
				req.response <- foodItem
			} else {
				req.available <- false
			}
		}

		time.Sleep(1100 * time.Millisecond) // Processing overhead (1.1 seconds = 1.1 hours simulation)
	}

	elapsed := time.Since(gk.startTime)
	fmt.Printf("\n\nHours of operation prior to closing: %.3f\n", elapsed.Seconds())

	gk.printSalesSummary() // HIGHLIGHT: added closing report, all grading options require it
}

// printSalesSummary prints per-FoodType and grand-total units sold. HIGHLIGHT: whole function is new
func (gk *GateKeeper) printSalesSummary() {
	fmt.Println("\n=================== Sales Summary ===================")
	fmt.Printf("%-12s | %s\n", "Food Type", "Units Sold")
	fmt.Println("------------------------------------------------------")
	for ft := food.Wheat; ft <= food.Fowel; ft++ {
		fmt.Printf("%-12s | %d\n", ft, gk.soldCounts[ft])
	}
	fmt.Println("------------------------------------------------------")
	fmt.Printf("%-12s | %d\n", "TOTAL", gk.totalSold)
	fmt.Println("=======================================================")
}

// HIGHLIGHT END
