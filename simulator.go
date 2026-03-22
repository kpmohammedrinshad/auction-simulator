package main

import (
	"fmt"
	"sync"
	"time"
)

// NewSimulator applies resource constraints and initialises TotalBidders bidders.
func NewSimulator(cfg ResourceConfig) *Simulator {
	cfg.Apply()

	bidders := make([]Bidder, TotalBidders)
	for i := range bidders {
		bidders[i] = Bidder{ID: i + 1}
	}
	return &Simulator{Bidders: bidders, Config: cfg}
}

// Run launches all TotalAuctions auctions concurrently, waits for every one to
// finish (or timeout), then prints per-auction results and the total elapsed time.
//
// Timing:
//   - globalStart  — captured just before the first goroutine is spawned.
//   - totalElapsed — measured only after wg.Wait() confirms all auctions are done.
func (s *Simulator) Run() {
	results := make([]AuctionResult, TotalAuctions)
	var wg sync.WaitGroup

	globalStart := time.Now() // ← start of the FIRST auction

	for i := 0; i < TotalAuctions; i++ {
		wg.Add(1)
		go func(auctionID int) {
			defer wg.Done()
			results[auctionID-1] = RunAuction(auctionID, s.Bidders)
		}(i + 1)
	}

	wg.Wait()
	totalElapsed := time.Since(globalStart) // ← completion of the LAST auction

	s.printResults(results, totalElapsed)
}

// printResults prints each auction's winner (or no-bid notice) and the overall timing.
func (s *Simulator) printResults(results []AuctionResult, total time.Duration) {
	fmt.Println("\n========== AUCTION RESULTS ==========")
	for _, r := range results {
		if r.WinnerID == -1 {
			fmt.Printf("Auction #%02d | No bids received | Bids: %d | Duration: %s\n",
				r.AuctionID, r.TotalBids, r.Duration.Round(time.Millisecond))
		} else {
			fmt.Printf("Auction #%02d | Winner: Bidder #%03d | WinBid: $%.2f | Bids: %d | Duration: %s\n",
				r.AuctionID, r.WinnerID, r.WinBid, r.TotalBids, r.Duration.Round(time.Millisecond))
		}
	}
	fmt.Println("=====================================")
	fmt.Printf("[Timing] First auction started → last auction completed: %s\n\n",
		total.Round(time.Millisecond))
}
