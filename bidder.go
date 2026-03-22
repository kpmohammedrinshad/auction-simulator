package main

import (
	"math/rand"
	"sync"
	"time"
)

// PlaceBid simulates a bidder's response to an auction item.
// It introduces a random delay (network / think time), then either skips
// (based on BidResponseRate) or calculates a bid as a weighted sum of item
// attributes and sends it on the shared channel.
func (b Bidder) PlaceBid(item Item, ch chan<- Bid, wg *sync.WaitGroup) {
	defer wg.Done()

	// Simulate variable network / decision latency
	time.Sleep(time.Duration(rand.Intn(BidderDelay)) * time.Millisecond)

	// Not every bidder responds — honour the response-rate probability
	if rand.Float64() > BidResponseRate {
		return
	}

	// Calculate bid: random weighted sum of all item attributes
	var amount float64
	for _, attr := range item.Attributes {
		amount += attr * (rand.Float64() + 0.5) // weight in range [0.5, 1.5)
	}

	ch <- Bid{BidderID: b.ID, Amount: amount}
}
