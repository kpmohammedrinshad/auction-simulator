package main

import (
	"sync"
	"time"
)

// RunAuction conducts a single auction identified by id.
//
// Flow:
//  1. Generate item with 20 attributes.
//  2. Launch all bidders as goroutines; each may send a bid on bidCh.
//  3. Collect bids until either all bidders finish OR the timeout fires.
//  4. Return an AuctionResult with the highest bid (or no-winner sentinel).
func RunAuction(id int, bidders []Bidder) AuctionResult {
	start := time.Now()
	item := NewItem(id)

	// Buffered channel prevents goroutine leaks when timeout fires early
	bidCh := make(chan Bid, TotalBidders)
	var wg sync.WaitGroup

	// Fan-out: every bidder evaluates the item concurrently
	for _, b := range bidders {
		wg.Add(1)
		go b.PlaceBid(item, bidCh, &wg)
	}

	// Close the channel once all bidders are done so the collect loop can exit
	go func() {
		wg.Wait()
		close(bidCh)
	}()

	// Collect bids; break on timeout or channel close (whichever comes first)
	timeout := time.After(AuctionTimeout)
	best := Bid{BidderID: -1} // -1 = no winner sentinel
	count := 0

collect:
	for {
		select {
		case bid, ok := <-bidCh:
			if !ok {
				break collect // channel closed — all bidders responded
			}
			count++
			if bid.Amount > best.Amount {
				best = bid // track highest bid
			}
		case <-timeout:
			break collect // auction time limit reached
		}
	}

	return AuctionResult{
		AuctionID: id,
		WinnerID:  best.BidderID,
		WinBid:    best.Amount,
		TotalBids: count,
		Duration:  time.Since(start),
	}
}
