// The software suite consists of main.go,
// food, stats, gatekeeper,
// producer, foodstore, and sales packages.
//
// This represents the software to manage an "embedded" planetary system
// food receiving and distribution system.

package main

import (
	"fmt"
	"time"

	"fooddistro/gatekeeper"
	"fooddistro/producer"
	"fooddistro/sales"
	"fooddistro/stats"
)

// main reads the requested concurrency settings and starts the simulation.
func main() {
	var numProductGenerators int
	var numPOS int

	fmt.Print("How many Product Generators? ")
	fmt.Scan(&numProductGenerators)

	fmt.Print("How many points of sale? ")
	fmt.Scan(&numPOS)
	fmt.Println()

	// Initialize components
	statistics := stats.NewStats()
	gk := gatekeeper.NewGateKeeper(20) // Default capacity of 20

	// Create producers
	for i := 0; i < numProductGenerators; i++ {
		pg := producer.NewProductGenerator(i+1, gk, statistics)
		pg.Start()
	}

	// Create sales points
	for i := 0; i < numPOS; i++ {
		rs := sales.NewRetailSales(i+1, gk, statistics)
		rs.Start()
	}

	// Block until system terminates (handled by gatekeeper timeout or rejection limit)
	time.Sleep(45 * time.Second) // Extra 5 seconds buffer for cleanup
	fmt.Println("\nSimulation completed.")
}
