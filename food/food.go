// In file food/food.go
//
// The software suite consists of main.go,
// food, stats, gatekeeper,
// producer, foodstore, and sales packages.
//
// This represents the software to manage an "embedded" planetary system
// food receiving and distribution system.

package food

import (
	"math/rand"
	"time"
)

type FoodType int

const (
	Wheat FoodType = iota
	Beans
	Corn
	Rice
	Potatoes
	Squash
	Tomato
	Steak
	Pork
	Fish
	Fowel
)

// String returns the display name for a food type.
func (ft FoodType) String() string {
	return []string{"Wheat", "Beans", "Corn", "Rice", "Potatoes", "Squash", "Tomato", "Steak", "Pork", "Fish", "Fowel"}[ft]
}

type FoodPack struct {
	FoodType     FoodType
	FoodShipment byte
}

// NewFoodPack creates a FoodPack with its type and shipment category.
func NewFoodPack(foodType FoodType, shipment byte) FoodPack {
	return FoodPack{
		FoodType:     foodType,
		FoodShipment: shipment,
	}
}

// init seeds the legacy package-level random number generator.

// init seeds the legacy package-level random number generator.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomFoodType returns a randomly selected food type.
// RandomFoodType returns a randomly selected food type.
func RandomFoodType() FoodType {
	return FoodType(rand.Intn(int(Fowel) + 1))
}

// IsGrainVegetable reports whether ft is a grain or vegetable type.
// IsGrainVegetable reports whether ft is a grain or vegetable type.
func IsGrainVegetable(ft FoodType) bool {
	return ft >= Wheat && ft <= Tomato
}
