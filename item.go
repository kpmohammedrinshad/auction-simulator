package main

import "math/rand"

// NewItem creates an Item with the given ID and TotalAttributes random float values.
// Each attribute is in the range [0, 1000) and represents one property of the object.
func NewItem(id int) Item {
	item := Item{ID: id}
	for i := range item.Attributes {
		item.Attributes[i] = rand.Float64() * 1000
	}
	return item
}
