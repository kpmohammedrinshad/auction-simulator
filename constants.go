package main

import "time"

const (
	TotalBidders    = 100             // Total number of bidders participating
	TotalAttributes = 20              // Number of attributes describing each auction item
	TotalAuctions   = 40              // Number of auctions running concurrently
	AuctionTimeout  = 3 * time.Second // Maximum wait time before an auction closes
	BidderDelay     = 500             // Max random delay (ms) simulating bidder think-time
	BidResponseRate = 0.75            // Probability (0–1) that a bidder actually places a bid
)
